package service

import (
	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) CreateSegment(slug string) (id int, affected bool, err error) {
	return s.repository.CreateSegment(slug)
}

func (s *Service) DeleteSegment(slug string) (deleted bool, err error) {
	return s.repository.DeleteSegment(slug)
}

func (s *Service) AddUserToSegments(userId int, slugs []string) error {
	segments, err := s.repository.GetSegments(slugs)
	if err != nil {
		return err
	}

	return s.repository.AddUserToSegments(userId, segments)
}

func (s *Service) RemoveUserFromSegments(userId int, slugs []string) error {
	segments, err := s.repository.GetSegments(slugs)
	if err != nil {
		return err
	}

	return s.repository.RemoveUserFromSegments(userId, segments)
}
