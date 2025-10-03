package main

import (
	"log"

	"github.com/gabrielmrtt/taski/cmd/seed/seeders"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	roledatabase "github.com/gabrielmrtt/taski/internal/role/infra/database"
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
	permissionSeeder := seeders.NewPermissionSeeder(roledatabase.NewPermissionBunRepository(coredatabase.DB))

	runSeeder(permissionSeeder.Run, "permissions")

	rolesSeeder := seeders.NewRolesSeeder(roledatabase.NewRoleBunRepository(coredatabase.DB), roledatabase.NewPermissionBunRepository(coredatabase.DB))

	runSeeder(rolesSeeder.Run, "roles")
}
