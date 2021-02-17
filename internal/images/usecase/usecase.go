package usecase

import (
	"github.com/Grishameister/Coursach/internal/images"
	"time"
)

type FrameUsecase struct {
	repo images.IRepository
}

func NewUsecase(repo images.IRepository) *FrameUsecase {
	return &FrameUsecase{
		repo: repo,
	}
}

func (uc *FrameUsecase) PostFrame(bytes []byte) error {
	return uc.repo.PostFrame(bytes)
}

func (uc *FrameUsecase) GetFrameByDate(date time.Time) []byte {
	return uc.repo.GetFrameByDate(date)
}

func (uc *FrameUsecase) GetLastFrame() []byte {
	return uc.repo.GetLastFrame()
}
