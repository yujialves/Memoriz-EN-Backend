package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"./controllers"
	"./driver"
)

func main() {

	db := driver.ConnectDB()
	controller := controllers.Controller{}

	router := mux.NewRouter()

	// エンドポイント
	router.Handle("/subjects", controller.SubjectsHandler(db)).Methods("GET")
	router.Handle("/question", controller.QuestionHandler(db)).Methods("POST")
	router.Handle("/question/correct", controller.CorrectHandler(db)).Methods("POST")
	router.Handle("/question/incorrect", controller.InCorrectHandler(db)).Methods("POST")

	// サーバーの起動
	log.Fatal(http.ListenAndServe(":9000", router))
	log.Println("サーバー起動 : 9000 port で受信")

}
