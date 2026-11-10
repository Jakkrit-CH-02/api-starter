package entity

import "gorm.io/gorm"

// User is the users table.
type User struct {
	gorm.Model
	FirstName string `gorm:"size:100"`
	LastName  string `gorm:"size:100"`
	// SQL Server cannot index nvarchar(MAX), so indexed columns need a size.
	Email    string `gorm:"size:255;uniqueIndex"`
	Password string `gorm:"size:255" json:"-"` // bcrypt hash

	// Nullable so that deleting a role sets it to NULL.
	UserRoleID *uint
	UserRole   *UserRole `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
