package main

import (
	"api-starter/pkg/db"
	"api-starter/pkg/env"
	"encoding/json"
	"fmt"

	externalApi "api-starter/external"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	app := fiber.New(fiber.Config{
		// BodyLimit: 100 * 1024 * 1024
		AppName: env.Env().MODULE_NAME,
		Immutable: true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(recover.New())
	app.Use(cors.New())

	db, _ := db.GetDB()
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	externalApi.PublicRoutes(app, db)

	err := app.Listen(fmt.Sprintf(":%v", env.Env().PORT))
	if err != nil {
		fmt.Printf("Error app listen: %v \n", err)
	}

}
