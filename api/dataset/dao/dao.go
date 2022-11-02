package dao

import (
	"gorm.io/gorm"
	"time"
)

const (
	StatusCreated = 1
	StatusDeleted = 2
)

type Dataset struct {
	gorm.Model
	ID       uint64    `gorm:"column:id,AUTO_INCREMENT"`
	Name     string    `gorm:"column:name,uniqueIndex"`
	Creator  string    `gorm:"column:creator"`
	CreateAt time.Time `gorm:"column:create_at"`
	Status   uint      `gorm:"column:status"`
}

type Tabler interface {
	TableName() string
}

func (Dataset) TableName() string {
	return "datasets"
}
