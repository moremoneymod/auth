package config

import (
	"errors"
	"net"
	"os"
)

const (
	grpcHostName = "GRPC_HOST"
	grpcPortName = "GRPC_PORT"
)

type GRPCConfig struct {
	host string
	port string
}

func NewGRPCConfig() (*GRPCConfig, error) {
	host := os.Getenv(grpcHostName)
	if len(host) == 0 {
		return nil, errors.New("grps host not found")
	}

	port := os.Getenv(grpcPortName)
	if len(port) == 0 {
		return nil, errors.New("grps port not found")
	}
	return &GRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg GRPCConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
