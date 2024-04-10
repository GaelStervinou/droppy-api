package provider

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"os"
)

func UseGoogleAuth() {
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

	goth.UseProviders(
		google.New(
			googleClientKey,
			googleSecret,
			googleRedirectUri,
			"email", "profile",
		),
	)
}
