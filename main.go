package main

import (
	"fmt"
	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"go-api/authentication"
	"go-api/internal/repositories"
	"go-api/internal/services/account"
	"go-api/internal/storage/postgres"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	//TODO migrer ça vers un fichier de configuration
	googleClientKey := os.Getenv("GOOGLE_AUTH_CLIENT_KEY")
	if "" == googleClientKey {
		panic("GOOGLE_AUTH_CLIENT_KEY is not set")
	}
	googleSecret := os.Getenv("GOOGLE_AUTH_SECRET")
	if "" == googleSecret {
		panic("GOOGLE_SECRET is not set")
	}
	googleRedirectUri := os.Getenv("GOOGLE_AUTH_REDIRECT_URI")
	if "" == googleRedirectUri {
		panic("GOOGLE_AUTH_REDIRECT_URI is not set")
	}

	postgres.AutoMigrate()

	var wg sync.WaitGroup
	repo := repositories.Setup(&wg)
	defer repo.Disconnect()

	goth.UseProviders(
		google.New(
			googleClientKey,
			googleSecret,
			googleRedirectUri,
			"email", "profile",
		),
	)

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

	p.Get("/users/me", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("me")
		res.Write([]byte("Hello, world!"))
	})

	fmt.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}
