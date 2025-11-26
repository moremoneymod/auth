package auth_v1

import (
	"context"

	desc "github.com/moremoneymod/auth/pkg/auth_v1"
)

type authService interface {
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, username string, password string) (string, error)
}

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService authService
}

func NewImplementation(service authService) *Implementation {
	return &Implementation{
		authService: service,
	}
}
