package persistence

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/config"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
)

var messageRepo *MessageRepo
var userRepo *UserRepo

func TestMain(m *testing.M) {
	db := getDatabase()
	messageRepo = NewMessageReopsitory(db)
	userRepo = NewUserRepository(db)
	res := m.Run()
	os.Exit(res)
}

func getDatabase() *io.SQLDatabase {
	cfg, err := config.LoadConfig(context.Background())
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	Sqlsetting := &config.SQLDBSetting{
		SqlDSN:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}

	db, err := io.NewDatabase(Sqlsetting)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return db

}
