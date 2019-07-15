package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	db         = newDB()
	errorDBNil = errors.New("db is nil")
)

func newDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("sql.Open error:", err)
	}
	return db
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, errorDBNil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.ExecContext(ctx, query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, errorDBNil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.QueryContext(ctx, query, args...)
}
