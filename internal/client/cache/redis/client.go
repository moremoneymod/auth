package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/moremoneymod/auth/internal/config"
)

type Client struct {
	pool   *redis.Pool
	config config.RedisConfig
}

func NewClient(pool *redis.Pool, config config.RedisConfig) *Client {
	return &Client{
		pool:   pool,
		config: config,
	}
}

func (c *Client) Get(ctx context.Context, key string) (any, error) {
	var value any
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		var errEx error
		value, errEx = conn.Do("GET", key)
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (c *Client) Set(ctx context.Context, key string, value any) error {
	var errEx error
	err := c.execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		errEx = conn.Send("SET", key, value)
		if errEx != nil {
			return errEx
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) execute(ctx context.Context, handler func(ctx context.Context, conn redis.Conn) error) error {
	conn, err := c.getConnect(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	err = handler(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getConnect(ctx context.Context) (redis.Conn, error) {
	getConnTimeoutCtx, cancel := context.WithTimeout(ctx, c.config.ConnectionTimeout())
	defer cancel()

	conn, err := c.pool.GetContext(getConnTimeoutCtx)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}
