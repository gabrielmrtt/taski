package main

import (
	"log"

	"github.com/gabrielmrtt/taski/cmd/seed/seeders"
	role_database_postgres "github.com/gabrielmrtt/taski/internal/role/database/postgres"
)

func runSeeder(seeder func() error, name string) {
	log.Println("Running seeder: ", name)
	err := seeder()
	if err != nil {
		log.Fatalf("Error running seeder: %v", err)
	}
	log.Println("Seeder completed: ", name)
}

func main() {
	permissionSeeder := seeders.NewPermissionSeeder(role_database_postgres.NewPermissionPostgresRepository())

	runSeeder(permissionSeeder.Run, "permissions")
}
