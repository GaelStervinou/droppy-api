package routes

import (
	"encoding/json"
	"go-api/pkg/jwt_helper"
	"go-api/pkg/services/account"
	"net/http"
)

func RefreshTokenHandler(res http.ResponseWriter, req *http.Request, acc account.AccountServiceIface) {
	tokenString := req.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(res, "no token provided", http.StatusUnauthorized)
		return
	}
	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		http.Error(res, "invalid token format", http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[7:]
	token, err := jwt_helper.VerifyToken(tokenString)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenInfo, err := acc.LoginFromRefreshToken(token.Raw)
	if err != nil {
		if err.Error() == "refresh token not found" {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(tokenInfo)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(payload)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
