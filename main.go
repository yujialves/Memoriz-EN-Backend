package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"

	"./auth"
	"./controllers"
	"./driver"
)

func init() {
	gotenv.Load()
}

func main() {

	db := driver.ConnectDB()
	subjectsController := controllers.SubjectsController{}
	questionController := controllers.QuestionController{}
	authController := controllers.AuthController{}

	router := mux.NewRouter()

	// エンドポイント
	router.Handle("/subjects", subjectsController.SubjectsHandler(db)).Methods("GET")
	router.Handle("/question", questionController.QuestionHandler(db)).Methods("POST")
	router.Handle("/question/correct", questionController.CorrectHandler(db)).Methods("POST")
	router.Handle("/question/incorrect", questionController.InCorrectHandler(db)).Methods("POST")
	router.Handle("/auth/login", authController.LoginHandler(db)).Methods("POST")
	router.Handle("/auth/refresh", auth.JwtMiddleware.Handler(authController.RefreshHandler(db))).Methods("GET")

	// サーバーの起動
	log.Fatal(http.ListenAndServe(":9000", router))
	log.Println("サーバー起動 : 9000 port で受信")

}
