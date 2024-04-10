package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
	"go-api/authentication"
	"go-api/authentication/provider"
	"go-api/internal/repositories"
	"go-api/internal/services/account"
	"go-api/internal/storage/postgres"
	"go-api/pkg/jwt_helper"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	provider.UseGoogleAuth()

	postgres.AutoMigrate()

	var wg sync.WaitGroup
	repo := repositories.Setup(&wg)
	defer repo.Disconnect()

	authentication.Init()

	p := pat.New()

	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("callback")
		authentication.GoogleAuthHandler(res, req, &account.AccountService{
			Repo: repo,
		})
	})
	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("auth")
		gothic.BeginAuthHandler(res, req)
	})

	p.Use(CurrentUserMiddleware)
	p.Get("/users/{id}", func(res http.ResponseWriter, req *http.Request) {
		username := context.Get(req, "username").(string)

		user, err := repo.UserRepository.GetByEmail(req.Context(), username)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		id := strings.TrimSpace(req.URL.Query().Get(":id"))

		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(res, "Invalid user ID", http.StatusBadRequest)
			return
		}

		requestedUser, err := repo.UserRepository.GetById(req.Context(), uint(userID))

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if requestedUser == nil {
			http.Error(res, "User not found", http.StatusNotFound)
			return
		}

		if user.GetID() != uint(userID) {
			http.Error(res, "You are not authorized to access this resource", http.StatusForbidden)
			return
		}

		payload, err := json.Marshal(requestedUser)
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
	})

	fmt.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}

func CurrentUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		token, err := jwt_helper.VerifyToken(tokenString)

		if err != nil {
			http.Error(w, "Failed to parse JWT token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username := claims["username"].(string)

			//TODO peut-être faire une requête pour récupérer le user et passer le user direct dans le context
			context.Set(r, "username", username)
		}
		next.ServeHTTP(w, r)
	})
}
