//go:generate go run db_migrate.go
package main

import "api-starter/pkg/db"

func main() {

	db.Migrate()
}