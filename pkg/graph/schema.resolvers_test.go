package graph

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/config"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/generated"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/loader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
	"go.uber.org/zap"
)

func TestQueryResolver_Users(t *testing.T) {
	t.Run("get all users", func(t *testing.T) {
		resolver := getQueryResolver()
		ctx := context.Background()
		users, err := resolver.Users(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
		assert.Len(t, users, 2)

		want := []*model.User{
			{ID: "2", Name: "Fuga"},
			{ID: "1", Name: "Hoge"},
		}
		if diff := cmp.Diff(want, users); len(diff) != 0 {
			t.Error("users mismatch (-want + got)", diff)
		}

	})
}

func getQueryResolver() generated.QueryResolver {
	return getResolver().Query()
}

func getResolver() *Resolver {
	mysqlDatabase := getDatabase()
	repositories, err := persistence.NewReopsitories(mysqlDatabase)
	if err != nil {
		log.Println("failed to new repositories", zap.Error(err))
	}

	loaders := loader.NewLoaders(repositories)
	resolver := &Resolver{
		Repo:    repositories,
		Loaders: loaders,
	}
	return resolver
}

func getDatabase() *io.SQLDatabase {
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	sqlSetting := &config.SQLDBSetting{
		SqlDSN:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}

	db, err := io.NewDatabase(sqlSetting)
	if err != nil {
		log.Println(err.Error())
	}
	return db
}
