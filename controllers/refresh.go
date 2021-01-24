package controllers

import (
	"net/http"
	"time"

	"memoriz-en/models"
	"memoriz-en/utils"
)

func (c AuthController) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// HTTPリクエストのヘッダーのトークンからユーザーID、トークンの期限を抽出
		claims := utils.GetMapClaimsFromRequest(r)
		user := claims["user"].(string)
		expireDate := time.Unix(int64(claims["exp"].(float64)), 0)

		// トークンの期限が切れてたら続行しない
		if time.Now().After(expireDate) {
			errorResponse := models.ErrorResponse{Message: "リフレッシュトークンの期限切れ。"}
			utils.ResponseWithError(w, http.StatusUnauthorized, errorResponse)
			return
		}

		token, refreshToken, exp := utils.GenerateToken(user)

		response := models.RefreshResponse{User: user, Password: "", Token: token, RefreshToken: refreshToken, ExpireDate: exp}
		utils.ResponseJSON(w, response)

	}
}
