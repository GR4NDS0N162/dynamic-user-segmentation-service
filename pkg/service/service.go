package service

import (
	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/model"
	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) CreateSegment(slug string) (model.Segment, error) {
	return s.repository.CreateSegment(slug)
}
