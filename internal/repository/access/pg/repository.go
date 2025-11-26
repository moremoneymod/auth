package pg

import "github.com/moremoneymod/auth/internal/client/pg"

type Repository struct {
	client pg.Client
}

func NewRepository(client pg.Client) *Repository {
	return &Repository{client: client}
}
