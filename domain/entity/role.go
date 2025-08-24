package entity

type UserRole struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
