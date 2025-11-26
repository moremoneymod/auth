package redis

import (
	"context"
	"errors"

	"github.com/moremoneymod/auth/internal/client/cache/redis"
	serv "github.com/moremoneymod/auth/internal/model"
)

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrInvalidValueType = errors.New("invalid value type")
)

type Repo struct {
	cl *redis.Client
}

func NewRepository(cl *redis.Client) *Repo {
	return &Repo{cl: cl}
}

func (repo *Repo) Get(ctx context.Context, key string) (*serv.AccessInfoCache, error) {
	value, err := repo.cl.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, ErrKeyNotFound
	}

	convertedValue, ok := value.(string)
	if !ok {
		return nil, ErrInvalidValueType
	}

	return &serv.AccessInfoCache{Role: convertedValue, EndpointAddress: key}, nil
}
