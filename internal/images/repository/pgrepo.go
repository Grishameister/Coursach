package repository

import (
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

func (r *Repository) PostFrame(bytes []byte) error {
	return r.db.PostFrame(bytes)
}
