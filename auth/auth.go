package auth

import (
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

var validationKeyGetter = func(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("SECRET_STRING")), nil
}

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: validationKeyGetter,
	SigningMethod:       jwt.SigningMethodHS256,
})
