package service

import "github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/repository"

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
