package repository

import (
	"errors"

	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSegment(slug string) (id int, affected bool, err error) {
	segment := model.Segment{Slug: slug}
	result := r.db.Unscoped().FirstOrCreate(&segment, segment)

	id = segment.ID
	affected = result.RowsAffected != 0
	err = result.Error
	if err != nil || (!affected && segment.IsDel == 0) {
		return
	}

	segment.IsDel = 0
	result = r.db.Save(&segment)

	id = segment.ID
	affected = result.RowsAffected != 0
	err = result.Error
	return
}

func (r *Repository) DeleteSegment(slug string) (deleted bool, err error) {
	segment := model.Segment{Slug: slug}
	result := r.db.First(&segment, segment)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if result.Error != nil {
		return false, result.Error
	}

	result = r.db.Delete(&segment)
	return result.RowsAffected != 0, result.Error
}

func (r *Repository) GetSegments(slugs []string) ([]model.Segment, error) {
	var segments []model.Segment
	if len(slugs) == 0 {
		return segments, nil
	}

	result := r.db.Where("slug IN ?", slugs).Find(&segments)
	return segments, result.Error
}

func (r *Repository) AddUserToSegments(userId int, segments []model.Segment) error {
	user := model.User{ID: userId}
	result := r.db.FirstOrCreate(&user)
	if result.Error != nil {
		return result.Error
	}

	for _, segment := range segments {
		action := model.Action{UserID: userId, SegmentID: segment.ID}
		result = r.db.Last(&action, action)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) || action.Type == 1 { // Пользователь не сегментировался или был удален из этого сегмента
			action = model.Action{UserID: userId, SegmentID: segment.ID, Type: 0}
			result = r.db.Create(&action)
			if result.Error != nil {
				return result.Error
			}
		} else if result.Error != nil {
			return result.Error
		}
		// Пользователь в этом сегменте
	}

	return nil
}

func (r *Repository) RemoveUserFromSegments(userId int, segments []model.Segment) error {
	user := model.User{ID: userId}
	result := r.db.FirstOrCreate(&user)
	if result.Error != nil {
		return result.Error
	}

	for _, segment := range segments {
		action := model.Action{UserID: userId, SegmentID: segment.ID}
		result = r.db.Last(&action, action)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) || action.Type == 1 { // Пользователь не сегментировался или вне этого сегмента
			continue
		} else if result.Error != nil {
			return result.Error
		}

		// Пользователь был добавлен в этот сегмент, но не удален из него
		action = model.Action{UserID: userId, SegmentID: segment.ID, Type: 1}
		result = r.db.Create(&action)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (r *Repository) GetActiveSegments(userId int) (segments []model.Segment, err error) {
	user := model.User{ID: userId}
	err = r.db.FirstOrCreate(&user).Error
	if err != nil {
		return
	}

	var segmentIds []int
	err = r.db.Model(&model.Action{}).
		Select("segment_id").
		Where("user_id = ?", userId).
		Group("segment_id").
		Having("count(id) % 2 != 0").
		Pluck("segment_id", &segmentIds).Error
	if err != nil {
		return
	}

	if len(segmentIds) != 0 {
		err = r.db.Find(&segments, segmentIds).Error
	}
	return
}

func (r *Repository) GetActions(year int, month int) (actions []model.Action, err error) {
	err = r.db.Model(&model.Action{}).
		Unscoped().
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Where("EXTRACT(MONTH FROM created_at) = ?", month).
		Preload("Segment").
		Order("id").
		Find(&actions).Error
	return
}
