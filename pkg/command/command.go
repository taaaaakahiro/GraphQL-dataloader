package command

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/config"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/generated"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/handler"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/loader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/server"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/version"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	exitOK    = 0
	exitError = 1
)

func Run() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	// init logger

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "falied to setlup logger: %s\n", err)
		return exitError
	}
	defer logger.Sync()
	logger = logger.With(zap.String("version", version.Version))

	// load config
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		logger.Error("failed to load config", zap.Error(err))
		return exitError
	}

	// init listener
	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		logger.Error("failed to listen port", zap.Error(err))
		return exitError
	}
	logger.Info("server start .listeninfg", zap.Int("port", cfg.Port))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// init mysql
	sqlSetting := &config.SQLDBSetting{
		SqlDSN:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}
	mysqlDatabase, err := io.NewDatabase(sqlSetting)
	if err != nil {
		logger.Error("failed to create mysql db repository", zap.Error(err), zap.String("DSN", cfg.DB.DSN))
		return exitError
	}

	repositories, err := persistence.NewReopsitories(mysqlDatabase)
	if err != nil {
		logger.Error("failed to new repositories", zap.Error(err))
		return exitError
	}

	// init loader
	loaders := loader.NewLoaders(repositories)
	query := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Repo:    repositories,
				Loaders: loaders,
			},
		}))

	// init to start GraphQL server
	registry := handler.NewHandler(logger, repositories, query, version.Version)
	httpServer := server.NewServer(registry, &server.Config{Log: logger})
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return httpServer.Serve(listener)
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	if err := httpServer.GracefulShutdown(); err != nil {
		return exitError
	}
	cancel()
	if err := wg.Wait(); err != nil {
		return exitError
	}

	return exitOK
}
