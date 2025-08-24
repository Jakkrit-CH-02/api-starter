package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName  string
	LastName   string
	Email      string `gorm:"uniqueIndex"`
	Password   string
	UserRoleID UserRole `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
