package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/moremoneymod/auth/internal/api/access_v1"
	"github.com/moremoneymod/auth/internal/api/auth_v1"
	"github.com/moremoneymod/auth/internal/client/pg"
	"github.com/moremoneymod/auth/internal/closer"
	"github.com/moremoneymod/auth/internal/config"
	"github.com/moremoneymod/auth/internal/model"
	access2 "github.com/moremoneymod/auth/internal/repository/access/pg"
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

type GRPCConfig interface {
	Address() string
}

type ServiceProvider struct {
	pgConfig   *config.PGConfig
	grpcConfig *config.GRPCConfig
	authConfig *config.AuthConfig

	pgClient         pg.Client
	userRepository   UserRepository
	accessRepository AccessRepository

	authService   AuthService
	accessService AccessService

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

func (s *ServiceProvider) AuthService(ctx context.Context) AuthService {
	if s.authService == nil {
		s.authService = auth.NewService(s.GetUserRepository(ctx), s.GetAuthConfig())
	}

	return s.authService
}

func (s *ServiceProvider) AccessService(ctx context.Context) AccessService {
	if s.accessService == nil {
		s.accessService = access.NewService(s.GetAccessRepository(ctx))
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
