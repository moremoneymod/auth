package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/moremoneymod/auth/internal/model"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func GenerateToken(userInfo *model.User, secretKey []byte, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userInfo.Username
	claims["role"] = userInfo.Role
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return secretKey, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userClaims := &model.UserClaims{
		Username: claims["username"].(string),
		Role:     claims["role"].(string),
	}

	return userClaims, nil
}
