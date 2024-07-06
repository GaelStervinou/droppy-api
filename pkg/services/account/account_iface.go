package account

import "context"

type AccountServiceIface interface {
	Create(string, string, string) error
	CreateWithGoogle(string, string) error
	Login(string, string) (*TokenInfo, error)
	LoginWithFirebase(string, context.Context) (*TokenInfo, error)
	LoginWithGoogle(string) (*TokenInfo, error)
	Logout(uint) error
	LoginFromRefreshToken(string) (*TokenInfo, error)
	EmailExists(string) (bool, error)
}

type TokenInfo struct {
	JWTToken     string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
	Expiry       int
}
