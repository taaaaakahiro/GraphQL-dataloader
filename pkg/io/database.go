package io

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	errs "github.com/pkg/errors"
)

type MySQLSettings interface {
	DSN() string
	MaxOpenConns() int
	MaxIdleConns() int
	ConnsMaxLifetime() int
}

type SQLDatabase struct {
	database *sql.DB
}

func NewDatabase(setting MySQLSettings) (*SQLDatabase, error) {
	db, err := sql.Open("mysql", setting.DSN())
	if err != nil {
		return nil, errs.WithStack(err)
	}

	// Check config
	if setting.MaxOpenConns() <= 0 {
		return nil, errs.WithStack(errs.New("required set max open conns"))
	}

	if setting.MaxIdleConns() <= 0 {
		return nil, errs.WithStack(errs.New("required set max idle conns"))
	}
	if setting.ConnsMaxLifetime() <= 0 {
		return nil, errs.WithStack(errs.New("required set conns max lifetime"))
	}
	db.SetMaxOpenConns(setting.MaxOpenConns())
	db.SetMaxIdleConns(setting.MaxIdleConns())
	db.SetConnMaxLifetime(time.Duration(setting.ConnsMaxLifetime()) * time.Second)

	return &SQLDatabase{database: db}, nil
}

func (d *SQLDatabase) Prepare(query string) (*sql.Stmt, error) {
	if d.database == nil {
		return nil, errDoesNotDB()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt, err := d.database.PrepareContext(ctx, query)
	return stmt, err
}

func errDoesNotDB() error {
	return errs.New("database does not exist. Please Open() first")
}
