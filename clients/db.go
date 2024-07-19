package clients

import (
	"database/sql"
	"log"

	"loan-payment/configs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var dbClient *sqlx.DB

func getDatabase() *sqlx.DB {
	if dbClient != nil {
		return dbClient
	}

	conf := configs.Get()

	if conf.DBMaster == nil {
		log.Fatalf("failed to get DB config")
	}

	db, err := sqlx.Open("mysql", conf.DBMaster.ConnectionString)
	if err != nil {
		log.Fatalf("failed to open DB master connection. %+v", err)
	}
	db.SetMaxIdleConns(conf.DBMaster.MaxIdle)
	db.SetMaxOpenConns(conf.DBMaster.MaxOpen)
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping DB master. %+v", err)
	}

	dbClient = db
	return dbClient
}

func DBBeginTransaction() (*sql.Tx, error) {
	var err error
	if tx, err := getDatabase().DB.Begin(); err == nil {
		return tx, nil
	}
	return nil, errors.Wrap(err, "failed to start tx")
}

func DBRollbackTransaction(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil {
		return errors.Wrap(err, "failed to rollback tx")
	}
	return nil
}

func DBCommitTransaction(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}
	return nil
}
