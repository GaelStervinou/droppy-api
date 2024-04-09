package account

import "context"

// TODO peut-être créer des objets pour les paramètres de création et de login
type AccountServiceIface interface {
	Create(context.Context, string, string, string, string) error
	CreateWithGoogle(context.Context, string, string, string, string) error
	Login(context.Context, string, string) (*TokenInfo, error)
	LoginWithGoogle(context.Context, string) (*TokenInfo, error)
	Logout(context.Context, string) error
	EmailExists(context.Context, string) (bool, error)
}

type TokenInfo struct {
	Token  string
	Expiry int
}
