package repository

import (
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
