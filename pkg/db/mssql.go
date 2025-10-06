package db

import (
	"api-starter/domain/entity"
	"api-starter/pkg/env"
	"context"
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func GetDB() (*gorm.DB, error) {
	server := env.Env().DBServer
	port := env.Env().DBPort
	dbName := env.Env().DBName
	fedauth := env.Env().FedAuth
	user := env.Env().UserName
	password := env.Env().Password

	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;fedauth=%s", server, user, password, port, dbName, fedauth)
	dial := sqlserver.New(sqlserver.Config{DriverName: "sqlserver", DSN: dsn})
	db, _ := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DryRun: false,
	})

	return db, nil
}

func Migrate() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery, // generate mode
	})

	db, _ := GetDB()
	g.UseDB(db) // reuse your gorm db

	// db.Migrator().DropTable(&entity.User{})

	db.AutoMigrate(
		&entity.User{},
		&entity.UserRole{},
	)

	g.ApplyInterface(func(Querier) { context.TODO() },
		entity.User{},
		entity.UserRole{},
	)

	// Generate the code
	g.Execute()
}
