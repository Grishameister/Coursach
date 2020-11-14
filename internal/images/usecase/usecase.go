package usecase

import (
	"github.com/Grishameister/Coursach/internal/images"
)

type FrameUsecase struct {
	repo images.IRepository
}

func NewUsecase(repo images.IRepository) *FrameUsecase {
	return &FrameUsecase{
		repo: repo,
	}
}

func (uc *FrameUsecase) PostFrame(bytes []byte) error{
	return uc.repo.PostFrame(bytes)
}
