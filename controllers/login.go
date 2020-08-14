package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"../models"
	"../utils"
)

type AuthController struct{}

func (c AuthController) LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// POSTの解析
		var post models.LoginPost
		json.NewDecoder(r.Body).Decode(&post)

		// バリデーション
		exp := regexp.MustCompile(`^[a-zA-Z][\w]{7,}$`)
		if !exp.MatchString(post.User) {
			errorResponse := models.ErrorResponse{Message: "ログインに失敗しました。"}
			utils.ResponseWithError(w, http.StatusBadRequest, errorResponse)
		}
		if len(post.Password) < 8 {
			errorResponse := models.ErrorResponse{Message: "ログインに失敗しました。"}
			utils.ResponseWithError(w, http.StatusBadRequest, errorResponse)
		}

		var hashedPassword string
		row := db.QueryRow("SELECT password FROM users WHERE user = ?", post.User)
		err := row.Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				errorResponse := models.ErrorResponse{Message: "ログインに失敗しました。"}
				utils.ResponseWithError(w, http.StatusBadRequest, errorResponse)
				return
			} else {
				log.Fatal(err)
			}
		}

		isValidPassword := utils.ComparePasswords(hashedPassword, []byte(post.Password))
		if isValidPassword {

			token, refreshToken, exp := utils.GenerateToken(post.User)
			if err != nil {
				log.Fatal(err)
			}

			response := models.LoginResponse{User: post.User, Password: "", Token: token, RefreshToken: refreshToken, ExpireDate: exp}
			utils.ResponseJSON(w, response)
		} else {
			errorResponse := models.ErrorResponse{Message: "ログインに失敗しました。"}
			utils.ResponseWithError(w, http.StatusUnauthorized, errorResponse)
		}
	}
}
