package loader

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/config"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
	"go.uber.org/zap"
)

var testLoaders *Loaders

func TestMain(m *testing.M) {
	db := getDatabase()
	repositories, err := persistence.NewReopsitories(db)
	if err != nil {
		log.Println("failed to new repositories", zap.Error(err))
		os.Exit(1)
	}
	testLoaders = NewLoaders(repositories)
	res := m.Run()
	os.Exit(res)
}

func TestNewLoaders(t *testing.T) {
	assert.NotNil(t, testLoaders.UserLoader)

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
		os.Exit(1)
	}
	return db
}
