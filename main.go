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
	router.Handle("/auth/login", authController.LoginHandler(db)).Methods("POST")
	router.Handle("/auth/refresh", auth.JwtMiddleware.Handler(authController.RefreshHandler())).Methods("GET")
	router.Handle("/subjects", auth.JwtMiddleware.Handler(subjectsController.SubjectsHandler(db))).Methods("GET")
	router.Handle("/question", auth.JwtMiddleware.Handler(questionController.QuestionHandler(db))).Methods("POST")
	router.Handle("/question/correct", auth.JwtMiddleware.Handler(questionController.CorrectHandler(db))).Methods("POST")
	router.Handle("/question/incorrect", auth.JwtMiddleware.Handler(questionController.InCorrectHandler(db))).Methods("POST")
	router.Handle("/question/list", auth.JwtMiddleware.Handler(questionController.QuestionListHandler(db))).Methods("POST")

	// サーバーの起動
	log.Fatal(http.ListenAndServe(":9000", router))
	log.Println("サーバー起動 : 9000 port で受信")

}
