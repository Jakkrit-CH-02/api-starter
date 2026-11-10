// Command gen migrates the database schema from the entities and generates
// the type-safe query code (gorm gen) into ./query.
//
// Run it from the project root:
//
//	go run ./gen
package main

import (
	"api-starter/domain/entity"
	"api-starter/pkg/db"

	"github.com/rs/zerolog/log"
	"gorm.io/gen"
)

func main() {
	gormDB, err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}
	defer db.Close(gormDB)

	// Referenced tables first: users has a foreign key to user_roles.
	if err := gormDB.AutoMigrate(
		&entity.UserRole{},
		&entity.User{},
	); err != nil {
		log.Error().Err(err).Msg("auto migrate failed")
		return
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery,
	})
	g.UseDB(gormDB)
	g.ApplyBasic(
		entity.UserRole{},
		entity.User{},
	)
	g.Execute()
}
