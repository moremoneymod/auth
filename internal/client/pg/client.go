package pg

import (
	"context"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}
type NamedExecer interface {
	GetContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

type PG interface {
	QueryExecer
	NamedExecer
	Pinger
	Close() error
}

type Client interface {
	Close() error
	PG() PG
}

type client struct {
	pg PG
}

func NewClient(ctx context.Context, pgCfg *pgxpool.Config) (*client, error) {
	dbc, err := pgxpool.ConnectConfig(ctx, pgCfg)
	if err != nil {
		log.Fatalf("failed to get db connection: %s", err.Error())
	}

	return &client{pg: &pg{pgxPool: dbc}}, nil
}

func (c *client) PG() PG {
	return c.pg
}

func (c *client) Close() error {
	if c.pg != nil {
		return c.pg.Close()
	}

	return nil
}
