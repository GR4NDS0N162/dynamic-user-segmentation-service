package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Action struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	Type      byte // 0 - add; 1 - remove
	UserID    int
	SegmentID int
	Segment   Segment
}

type User struct {
	ID      int `gorm:"primaryKey;autoIncrement:false"`
	Actions []Action
}

type Segment struct {
	ID      int                   `gorm:"primaryKey"`
	Slug    string                `gorm:"unique"`
	IsDel   soft_delete.DeletedAt `gorm:"softDelete:flag"`
	Actions []Action
}
