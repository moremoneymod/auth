package auth

import (
	"context"
	"errors"

	"github.com/moremoneymod/auth/internal/repository"
	"github.com/moremoneymod/auth/internal/service"
	"github.com/moremoneymod/auth/internal/utils"
)

func (s *Service) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, s.authConfig.RefreshTokenSecret())
	if errors.Is(err, utils.ErrInvalidToken) {
		return "", service.ErrInvalidToken
	}
	if err != nil {
		return "", err
	}

	userInfo, err := s.userRepository.Get(ctx, claims.Username)
	if errors.Is(err, repository.ErrUserNotFound) {
		return "", service.ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(userInfo, s.authConfig.AccessTokenSecret(), s.authConfig.AccessTokenExpiration())
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
