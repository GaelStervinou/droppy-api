package jwt_helper

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userId uint, roles []string) (string, string, int, error) {
	expiry := time.Now().Add(time.Minute * 5).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":   userId,
			"exp":   expiry,
			"roles": roles,
		})

	tokenString, err := jwtToken.SignedString(secretKey)

	if err != nil {
		return "", "", 0, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenExpiry := time.Now().AddDate(0, 0, 7).Unix()
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = userId
	rtClaims["exp"] = refreshTokenExpiry
	rt, err := refreshToken.SignedString(secretKey)

	if err != nil {
		return "", "", 0, err
	}

	return tokenString, rt, int(refreshTokenExpiry), nil
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

func GetUserIdFromToken(tokenString string) uint {
	token, _ := VerifyToken(tokenString)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := claims["sub"].(float64)
	return uint(userId)
}
