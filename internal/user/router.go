package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Router builds the user module (repository -> service -> handler) and
// registers its routes under /users.
func Router(r fiber.Router, db *gorm.DB) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	group := r.Group("/users")
	group.Post("/", handler.Create)
	group.Get("/", handler.List)
	group.Get("/:id", handler.GetByID)
}
