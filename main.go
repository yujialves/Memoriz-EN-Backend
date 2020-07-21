package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	// DB 接続
	db, err := sql.Open("mysql", USER+":"+PASSWORD+"@tcp("+HOSTNAME+")/"+DBNAME)
	if err != nil {
		log.Fatal(err)
	}

	// DB 疎通確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// DB 切断
	db.Close()

	router := mux.NewRouter()

	// エンドポイント

	// サーバーの起動
	log.Fatal(http.ListenAndServe(":9000", router))
	log.Println("サーバー起動 : 9000 port で受信")

}
