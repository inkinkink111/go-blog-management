package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Email  string `json:"email"`
	UserId string `json:"userId"`
	Exp    int64  `json:"exp"`
}

func GenerateToken(email string, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 8).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func VerifyToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return "", errors.New("could not parse token")
	}

	isTokenValid := parsedToken.Valid

	if !isTokenValid {
		return "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("invalid token claims")
	}

	userId := claims["userId"].(string)

	return userId, nil
}
