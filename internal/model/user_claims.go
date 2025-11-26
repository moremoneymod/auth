package model

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.Claims
	Username string `json:"username"`
	Role     string `json:"role"`
}
