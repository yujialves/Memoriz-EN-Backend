package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDB() *sql.DB {

	// DB 接続
	fmt.Println(os.Getenv("USER"))
	fmt.Println(os.Getenv("PASSWORD"))
	fmt.Println(os.Getenv("HOSTNAME"))
	fmt.Println(os.Getenv("DBNAME"))
	db, err := sql.Open("mysql", os.Getenv("USER")+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("HOSTNAME")+")/"+os.Getenv("DBNAME"))
	if err != nil {
		// log.Fatal(err)
		log.Println(err)
	}

	// DB 疎通確認
	err = db.Ping()
	if err != nil {
		log.Println(err)
	}

	return db
}
