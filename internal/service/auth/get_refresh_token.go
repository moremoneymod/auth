package auth

import (
	"context"
	"errors"

	"github.com/moremoneymod/auth/internal/repository"
	"github.com/moremoneymod/auth/internal/service"
	"github.com/moremoneymod/auth/internal/utils"
)

func (s *Service) GetRefreshToken(ctx context.Context, username string, password string) (string, error) {
	userInfo, err := s.userRepository.Get(ctx, username)
	if errors.Is(err, repository.ErrUserNotFound) {
		return "", service.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	if !utils.VerifyPassword(userInfo.Password, password) {
		return "", service.ErrInvalidCredentials
	}

	refreshToken, err := utils.GenerateToken(userInfo, s.authConfig.RefreshTokenSecret(), s.authConfig.RefreshTokenExpiration())
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
