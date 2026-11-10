package entity

// UserRole is the user_roles table.
type UserRole struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50"`
}
