package jwt_helper

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userId uint) (string, int, error) {
	expiry := time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    expiry,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, int(expiry), nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
