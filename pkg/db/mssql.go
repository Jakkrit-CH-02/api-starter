// Package db opens the SQL Server connection used by the service.
package db

import (
	"fmt"

	"api-starter/pkg/env"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a gorm connection to SQL Server using the DB_* settings.
func Connect() (*gorm.DB, error) {
	cfg := env.Env()

	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s",
		cfg.DBServer, cfg.DBUser, cfg.DBPassword, cfg.DBPort, cfg.DBName)
	if cfg.DBFedAuth != "" {
		dsn += ";fedauth=" + cfg.DBFedAuth
	}

	// Log every SQL statement in DEV, only slow queries and errors elsewhere.
	logLevel := logger.Warn
	if cfg.Build == "DEV" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(
		sqlserver.New(sqlserver.Config{DriverName: "sqlserver", DSN: dsn}),
		&gorm.Config{Logger: logger.Default.LogMode(logLevel)},
	)
	if err != nil {
		return nil, fmt.Errorf("open sqlserver %s/%s: %w", cfg.DBServer, cfg.DBName, err)
	}
	return db, nil
}

// Close closes the underlying connection pool.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
