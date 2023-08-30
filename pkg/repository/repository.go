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
