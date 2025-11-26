package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/moremoneymod/auth/internal/client/pg"
	serv "github.com/moremoneymod/auth/internal/model"
	"github.com/moremoneymod/auth/internal/repository"
	"github.com/moremoneymod/auth/internal/repository/user/converter"
	repo "github.com/moremoneymod/auth/internal/repository/user/model"
)

type Repository struct {
	client pg.Client
}

func NewRepository(client pg.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) Get(ctx context.Context, username string) (*serv.User, error) {
	builder := sq.Select("username", "password", "role").
		PlaceholderFormat(sq.Dollar).From("users").
		Where(sq.Eq{"username": username}).
		Limit(1)
	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "users.Get",
		QueryRaw: query}

	var user repo.User

	err = r.client.PG().GetContext(ctx, &user, q, v...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}
