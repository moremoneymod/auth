package config

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

type RedisConfig struct {
	host string
	port string

	connectionTimeout time.Duration

	maxIdle     int
	idleTimeout time.Duration
}

func NewRedisConfig() (*RedisConfig, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return nil, errors.New("redis host not set")
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		return nil, errors.New("redis port not set")
	}

	maxIdleStr := os.Getenv("REDIS_MAX_IDLE")
	if maxIdleStr == "" {
		return nil, errors.New("redis max idle not set")
	}

	maxIdle, err := strconv.Atoi(maxIdleStr)
	if err != nil {
		return nil, err
	}

	idleTimeoutStr := os.Getenv("REDIS_MAX_IDLE_TIMEOUT")
	if idleTimeoutStr == "" {
		return nil, errors.New("redis idle timeout not set")
	}

	idleTimeout, err := time.ParseDuration(idleTimeoutStr)
	if err != nil {
		return nil, errors.New("redis idle timeout not set")
	}

	return &RedisConfig{
		host:        host,
		port:        port,
		maxIdle:     maxIdle,
		idleTimeout: idleTimeout,
	}, nil
}

func (rc *RedisConfig) Host() string {
	return rc.host
}

func (rc *RedisConfig) Port() string {
	return rc.port
}

func (rc *RedisConfig) Address() string {
	return net.JoinHostPort(rc.Host(), rc.Port())
}

func (rc *RedisConfig) MaxIdle() int {
	return rc.maxIdle
}

func (rc *RedisConfig) IdleTimeout() time.Duration {
	return rc.idleTimeout
}
