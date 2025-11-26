package app

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/moremoneymod/auth/internal/api/access_v1"
	"github.com/moremoneymod/auth/internal/api/auth_v1"
	redis2 "github.com/moremoneymod/auth/internal/client/cache/redis"
	"github.com/moremoneymod/auth/internal/client/pg"
	"github.com/moremoneymod/auth/internal/closer"
	"github.com/moremoneymod/auth/internal/config"
	"github.com/moremoneymod/auth/internal/model"
	access2 "github.com/moremoneymod/auth/internal/repository/access/pg"
	cacheRepo "github.com/moremoneymod/auth/internal/repository/access/redis"
	"github.com/moremoneymod/auth/internal/repository/user"
	"github.com/moremoneymod/auth/internal/service/access"
	"github.com/moremoneymod/auth/internal/service/auth"
)

type AuthService interface {
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, username string, password string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpointAddress string) (bool, error)
}

type UserRepository interface {
	Get(ctx context.Context, username string) (*model.User, error)
}

type AccessRepository interface {
}

type CacheRepository interface {
	Get(ctx context.Context, key string) (*model.AccessInfoCache, error)
}

type GRPCConfig interface {
	Address() string
}

type ServiceProvider struct {
	pgConfig    *config.PGConfig
	grpcConfig  *config.GRPCConfig
	authConfig  *config.AuthConfig
	redisConfig *config.RedisConfig

	pgClient         pg.Client
	redisClient      *redis2.Client
	userRepository   UserRepository
	accessRepository AccessRepository

	authService     AuthService
	accessService   AccessService
	cacheRepository CacheRepository

	authImpl   *auth_v1.Implementation
	accessImpl *access_v1.Implementation
}

func newServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) GetPGConfig() *config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *ServiceProvider) GetRedisConfig() *config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %v", err)
		}

		s.redisConfig = cfg
	}
	return s.redisConfig
}

func (s *ServiceProvider) GetGRPCConfig() GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *ServiceProvider) GetAuthConfig() *config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config: %v", err)
		}
		s.authConfig = cfg
	}

	return s.authConfig
}

func (s *ServiceProvider) GetPGClient(ctx context.Context) pg.Client {
	if s.pgClient == nil {
		pgCfg, err := pgxpool.ParseConfig(s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to get db config: %v", err)
		}

		cl, err := pg.NewClient(ctx, pgCfg)
		if err != nil {
			log.Fatalf("failed to create pg client: %v", err)
		}
		err = cl.PG().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)

		s.pgClient = cl
	}

	return s.pgClient
}

func (s *ServiceProvider) GetRedisClient(ctx context.Context) *redis2.Client {
	if s.redisClient == nil {
		client := redis2.NewClient(&redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", s.redisConfig.Address())
			},
			MaxIdle:     s.redisConfig.MaxIdle(),
			IdleTimeout: s.redisConfig.IdleTimeout(),
		}, s.GetRedisConfig())

		s.redisClient = client

	}
	return s.redisClient
}

func (s *ServiceProvider) GetUserRepository(ctx context.Context) UserRepository {
	if s.userRepository == nil {
		s.userRepository = user.NewRepository(s.GetPGClient(ctx))
	}

	return s.userRepository
}

func (s *ServiceProvider) GetAccessRepository(ctx context.Context) AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = access2.NewRepository(s.GetPGClient(ctx))
	}

	return s.accessRepository
}

func (s *ServiceProvider) GetCacheRepository(ctx context.Context) CacheRepository {
	if s.cacheRepository == nil {
		s.cacheRepository = cacheRepo.NewRepository(s.GetRedisClient(ctx))
	}
	return s.cacheRepository
}

func (s *ServiceProvider) AuthService(ctx context.Context) AuthService {
	if s.authService == nil {
		s.authService = auth.NewService(s.GetUserRepository(ctx), s.GetAuthConfig())
	}

	return s.authService
}

func (s *ServiceProvider) AccessService(ctx context.Context) AccessService {
	if s.accessService == nil {
		s.accessService = access.NewService(s.GetAccessRepository(ctx), s.GetCacheRepository(ctx), s.GetAuthConfig())
	}

	return s.accessService
}

func (s *ServiceProvider) AuthImpl(ctx context.Context) *auth_v1.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth_v1.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *ServiceProvider) AccessImpl(ctx context.Context) *access_v1.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access_v1.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}
