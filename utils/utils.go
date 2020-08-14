package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"../models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func ResponseWithError(w http.ResponseWriter, status int, error models.ErrorResponse) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ComparePasswords(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GenerateToken(user string) (token string, refreshToken string, expireDate int64) {

	secret := os.Getenv("SECRET_STRING")
	exp := getUnixMillis(time.Now().AddDate(0, 0, 7))

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "Memoriz-EN",
		"user": user,
		"exp":  exp,
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	claims = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "Memoriz-EN",
		"user": user,
		"exp":  getUnixMillis(time.Now().AddDate(0, 6, 0)),
	})

	refreshToken, err = claims.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	return token, refreshToken, exp
}

func getUnixMillis(exp time.Time) int64 {
	nanos := exp.UnixNano()
	millis := nanos / 1000000
	return millis
}
