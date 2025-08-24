package user

import (
	"gorm.io/gorm"
)

type (
	Repository interface {
	}
	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) repository {
	// query.SetDefault(db)
	return repository{db}
}
