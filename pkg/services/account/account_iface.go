package account

// TODO peut-être créer des objets pour les paramètres de création et de login
type AccountServiceIface interface {
	Create(string, string, string, string, string) error
	CreateWithGoogle(string, string, string, string) error
	Login(string, string) (*TokenInfo, error)
	LoginWithGoogle(string) (*TokenInfo, error)
	Logout(uint) error
	LoginFromRefreshToken(string) (*TokenInfo, error)
	EmailExists(string) (bool, error)
}

type TokenInfo struct {
	JWTToken     string
	RefreshToken string
	Expiry       int
}
