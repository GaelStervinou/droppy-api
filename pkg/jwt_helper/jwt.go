package jwt_helper

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("secret-key")

func GenerateToken(username string) (string, int, error) {
	expiry := time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      expiry,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, int(expiry), nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
