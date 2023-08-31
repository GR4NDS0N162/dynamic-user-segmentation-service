package service

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

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

func (s *Service) GetActiveSegments(userId int) (slugs []string, err error) {
	segments, err := s.repository.GetActiveSegments(userId)
	if err != nil {
		return
	}

	slugs = []string{}
	for _, segment := range segments {
		slugs = append(slugs, segment.Slug)
	}
	return
}

func (s *Service) GetFile(year int, month int) (filename string, err error) {
	actions, err := s.repository.GetActions(year, month)
	if err != nil {
		return
	}

	var lastId = 0
	if len(actions) != 0 {
		lastId = actions[len(actions)-1].ID
	}
	filename = fmt.Sprintf("%v-%v-%v.csv", year, month, lastId)

	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var records [][]string
	for _, action := range actions {
		var operation string
		if action.Type == 0 {
			operation = "Добавление"
		} else {
			operation = "Удаление"
		}
		records = append(records, []string{strconv.Itoa(action.UserID), action.Segment.Slug, operation, action.CreatedAt.String()})
	}

	err = writer.WriteAll(records)
	if err != nil {
		return
	}

	writer.Flush()
	err = writer.Error()

	return
}
