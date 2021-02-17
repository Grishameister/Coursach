package repository

import (
	"github.com/Grishameister/Coursach/internal/database"
	"time"
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

func (r *Repository) GetFrameByDate(date time.Time) []byte {
	return r.db.GetFrameByDate(date)
}

func (r *Repository) GetLastFrame() []byte {
	return r.db.GetLastFrame()
}
