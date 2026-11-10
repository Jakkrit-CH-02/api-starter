// Package external wires every HTTP route exposed by the service.
package external

import (
	"errors"

	"api-starter/internal/user"
	"api-starter/pkg/env"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PublicRoutes registers the swagger UI, the service info endpoint and the
// routers of every internal module.
func PublicRoutes(app *fiber.App, db *gorm.DB) {
	app.Get("/api-doc/*", swagger.HandlerDefault)

	v1 := app.Group("/api/v1")
	v1.Get("/", info)

	// Routers of internal modules. Add new modules here.
	user.Router(v1, db)
}

func info(c *fiber.Ctx) error {
	cfg := env.Env()
	return c.JSON(fiber.Map{
		"productCode": cfg.ProductCode,
		"moduleName":  cfg.ModuleName,
		"build":       cfg.Build,
		"release":     cfg.Release,
		"port":        cfg.Port,
	})
}

// ErrorHandler turns every error returned by a handler into a JSON response:
// {"message": "..."}.
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	return c.Status(code).JSON(fiber.Map{"message": message})
}
