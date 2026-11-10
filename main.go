package main

import (
	"os"
	"os/signal"
	"syscall"

	"api-starter/external"
	"api-starter/pkg/db"
	"api-starter/pkg/env"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := env.Env()

	gormDB, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

	app := fiber.New(fiber.Config{
		AppName:      cfg.ModuleName,
		Immutable:    true,
		ErrorHandler: external.ErrorHandler,
	})
	app.Use(recover.New())
	app.Use(cors.New())

	external.PublicRoutes(app, gormDB)

	// Stop accepting requests on Ctrl+C / SIGTERM so Listen returns and the
	// database connection below is closed cleanly.
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		log.Info().Msg("shutting down")
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("server shutdown")
		}
	}()

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Error().Err(err).Msg("server stopped")
	}

	if err := db.Close(gormDB); err != nil {
		log.Error().Err(err).Msg("close database")
	}
}
