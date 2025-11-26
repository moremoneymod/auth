package config

import (
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	accessTokenSecretKeyEnvName  = "ACCESS_TOKEN_SECRET_KEY"
	accessTokenExpirationEnvName = "ACCESS_TOKEN_EXPIRATION_MINUTES"

	refreshTokenSecretKeyEnvName  = "REFRESH_TOKEN_SECRET_KEY"
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION_MINUTES"
)

type AuthConfig struct {
	accessTokenSecretKey  []byte
	refreshTokenSecretKey []byte

	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

func NewAuthConfig() (*AuthConfig, error) {
	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)

	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)

	accessTokenExpiration, err := strconv.Atoi(os.Getenv(accessTokenExpirationEnvName))
	if err != nil {
		return nil, err
	}

	refreshTokenExpiration, err := strconv.Atoi(os.Getenv(refreshTokenExpirationEnvName))
	if err != nil {
		return nil, err
	}

	return &AuthConfig{
		accessTokenSecretKey:   []byte(accessTokenSecretKey),
		refreshTokenSecretKey:  []byte(refreshTokenSecretKey),
		accessTokenExpiration:  time.Minute * time.Duration(accessTokenExpiration),
		refreshTokenExpiration: time.Minute * time.Duration(refreshTokenExpiration),
	}, nil
}

func decode(key string) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key is empty")
	}

	decodeKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return decodeKey, nil
}

func (cfg *AuthConfig) AccessTokenSecret() []byte {
	return cfg.accessTokenSecretKey
}

func (cfg *AuthConfig) RefreshTokenSecret() []byte {
	return cfg.refreshTokenSecretKey
}

func (cfg *AuthConfig) AccessTokenExpiration() time.Duration {
	return cfg.accessTokenExpiration
}

func (cfg *AuthConfig) RefreshTokenExpiration() time.Duration {
	return cfg.refreshTokenExpiration
}
