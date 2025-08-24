package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Router(r fiber.Router, db *gorm.DB) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)


	groutRoute := r.Group("/users")

	groutRoute.Post("/")

}
