package repository

import (
	"context"
	"github.com/Grishameister/Coursach/internal/database"
)

type Repository struct {
	db database.DBInterface
}

func NewRepo(d database.DBInterface) *Repository {
	return &Repository{
		db: d,
	}
}

func (r *Repository) PostFrame(c context.Context, bytes []byte) error {
	return r.db.PostFrame(c, bytes)
}