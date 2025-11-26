package access

import (
	"context"

	"github.com/moremoneymod/auth/internal/config"
	"github.com/moremoneymod/auth/internal/model"
)

type accessRepository interface {
}

type accessCache interface {
	Get(ctx context.Context, key string) (*model.AccessInfoCache, error)
}

type Service struct {
	accessRepository accessRepository
	cacheRepository  accessCache
	authConfig       config.AuthConfig
}

func NewService(accessRepository accessRepository) *Service {
	return &Service{
		accessRepository: accessRepository,
	}
}
