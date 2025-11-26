package access

import (
	"context"
	"errors"
	"strings"

	"github.com/moremoneymod/auth/internal/service"
	"github.com/moremoneymod/auth/internal/utils"
	"google.golang.org/grpc/metadata"
)

const authPrefix = "Bearer "

func (s *Service) Check(ctx context.Context, endpointAddress string) (bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("no metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return false, errors.New("no authorization header")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return false, errors.New("invalid authorization header")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)
	claims, err := utils.VerifyToken(accessToken, s.authConfig.AccessTokenSecret())
	if err != nil {
		return false, service.ErrInvalidToken
	}

	role, err := s.cacheRepository.Get(ctx, endpointAddress)
	if err != nil {
		return true, err
	}

	if role.Role == claims.Role {
		return true, nil
	}

	return false, errors.New("access denied")
}
