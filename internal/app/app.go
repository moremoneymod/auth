package app

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/moremoneymod/auth/internal/closer"
	"github.com/moremoneymod/auth/internal/config"
	"github.com/moremoneymod/auth/internal/interceptor"
	"github.com/moremoneymod/auth/internal/metrics"
	"github.com/moremoneymod/auth/internal/rate_limiter"
	descAccess "github.com/moremoneymod/auth/pkg/access_v1"
	descAuth "github.com/moremoneymod/auth/pkg/auth_v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider  *ServiceProvider
	grpcServer       *grpc.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runPrometheusServer()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		metrics.Init,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initPrometheusServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {

	rateLimiter := rate_limiter.NewTokenBucketLimiter(ctx, 100, time.Second)

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.MetricsInterceptor,
				interceptor.NewRateLimiterInterceptor(rateLimiter).Unary)),
	)
	reflection.Register(a.grpcServer)
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))

	return nil
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) runGRPCServer() error {
	list, err := net.Listen("tcp", a.serviceProvider.GetGRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
