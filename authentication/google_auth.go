package authentication

import (
	"encoding/json"
	"github.com/markbates/goth/gothic"
	"go-api/pkg/services/account"
	"net/http"
)

func GoogleAuthHandler(res http.ResponseWriter, req *http.Request, acc account.AccountServiceIface) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		http.Error(res, "email is empty", http.StatusInternalServerError)
		return
	}

	isKnown, err := acc.EmailExists(req.Context(), user.Email)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if false == isKnown {
		err := acc.CreateWithGoogle(req.Context(), user.FirstName, user.LastName, user.Email, user.UserID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tokenInfo, err := acc.LoginWithGoogle(req.Context(), user.Email)

	if err != nil {
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
