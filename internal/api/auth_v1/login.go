package auth_v1

import (
	"context"

	desc "github.com/moremoneymod/auth/pkg/auth_v1"
)

func (s *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := s.authService.GetRefreshToken(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}
