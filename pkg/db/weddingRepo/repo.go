package weddingRepo

import (
	"github.com/babafemi99/WR/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r Repository) PersistWedding() error {
	panic("implement me ")
}

func (r Repository) ToggleWeddingLink(key string) error {
	panic("implement me ")
}

func (r Repository) GetLinkByKey(key string) (model.Wedding, error) {
	panic("implement me ")
}
