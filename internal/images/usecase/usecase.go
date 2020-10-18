package usecase

import (
	"context"
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

func (uc *FrameUsecase) PostFrame(c context.Context, bytes []byte) error{
	return uc.repo.PostFrame(c, bytes)
}