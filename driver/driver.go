package driver

import (
	"database/sql"
	"log"

	"../secret"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDB() *sql.DB {

	// DB 接続
	db, err := sql.Open("mysql", secret.DSN)
	if err != nil {
		log.Fatal(err)
	}

	// DB 疎通確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
