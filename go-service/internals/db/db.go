package db

import (
	"context"
	"database/sql"
	"go-servie/dbmodel"
	"go-servie/utils"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *dbmodel.Queries

func ConnectDB() (*sql.DB, error) {
	log.Println("Connecting to database")
	var dsn string
	if utils.APP_ENV == "dev" {
		dsn = utils.DEVELOPMENT_DATABASE_URL
	} else {
		dsn = utils.PRODUCTION_DB_URL
	}
	if dsn == "" {
		log.Fatal("DB_URL env variable not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("failed to ping db:", err)
		return nil, err
	}
	log.Println("Connected to database")
	DB = dbmodel.New(db)
	return db, nil
}
