package external

import (
	"api-starter/pkg/env"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

func PublicRoutes(app *fiber.App, db *gorm.DB) {
	routerV1 := app.Group("api/v1")
	routerSWG := app.Group("api-doc")

	routerSWG.Use("", swagger.HandlerDefault)
	routerV1.Get("", func(c *fiber.Ctx) error {
		return c.JSON(map[string]any{
			"ProductCode": env.Env().PRODUCT_CODE,
			"moduleName":  env.Env().MODULE_NAME,
			"Build":       env.Env().BUILD,
			"Release":     env.Env().RELEASE,
			"Port":        env.Env().PORT,
		})
	})

	// router for internal apis
	// example: users
	// UserRouter(routerV1, db)
}
