package repository

import (
	"fmt"

	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSegment(slug string) (model.Segment, error) {
	segment := model.Segment{Slug: slug}
	result := r.db.Unscoped().FirstOrCreate(&segment, segment)

	if result.Error != nil {
		return segment, result.Error
	}

	if result.RowsAffected == 0 && segment.IsDel == 0 {
		return segment, fmt.Errorf("segment %s already exists", slug)
	}

	if segment.IsDel == 1 {
		segment.IsDel = 0
		r.db.Save(&segment)
	}

	return segment, nil
}
