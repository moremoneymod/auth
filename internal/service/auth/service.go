package auth

import (
	"context"

	"github.com/moremoneymod/auth/internal/config"
	"github.com/moremoneymod/auth/internal/model"
)

type userRepository interface {
	Get(ctx context.Context, username string) (*model.User, error)
}

type Service struct {
	userRepository userRepository
	authConfig     config.AuthConfig
}

func NewService(userRepository userRepository, authConfig *config.AuthConfig) *Service {
	return &Service{userRepository: userRepository,
		authConfig: *authConfig}
}
