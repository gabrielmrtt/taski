package main

import (
	"log"
	"os"

	"github.com/gabrielmrtt/taski/cmd/seed/seeders"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	roledatabase "github.com/gabrielmrtt/taski/internal/role/infra/database"
	"github.com/uptrace/bun"
)

type Seeder interface {
	Name() string
	Run() error
}

func runSeeder(seeder func() error, name string) {
	log.Println("Running seeder: ", name)
	err := seeder()
	if err != nil {
		log.Fatalf("Error running seeder: %v", err)
	}
	log.Println("Seeder completed: ", name)
}

func main() {
	var env string = "default"
	if len(os.Args) >= 2 {
		env = os.Args[1]
	}

	var connection *bun.DB
	var seedersArr []Seeder = make([]Seeder, 0)

	switch env {
	case "default":
		connection = coredatabase.GetPostgresConnection()
		permissionRepository := roledatabase.NewPermissionBunRepository(connection)
		roleRepository := roledatabase.NewRoleBunRepository(connection)

		seedersArr = append(seedersArr, seeders.NewPermissionSeeder(seeders.PermissionSeederOptions{
			PermissionRepository: permissionRepository,
		}))

		seedersArr = append(seedersArr, seeders.NewRolesSeeder(seeders.RolesSeederOptions{
			RoleRepository:       roleRepository,
			PermissionRepository: permissionRepository,
		}))
	case "test":
		connection = coredatabase.GetSQLiteConnection()
		permissionRepository := roledatabase.NewPermissionBunRepository(connection)
		roleRepository := roledatabase.NewRoleBunRepository(connection)

		seedersArr = append(seedersArr, seeders.NewPermissionSeeder(seeders.PermissionSeederOptions{
			PermissionRepository: permissionRepository,
		}))

		seedersArr = append(seedersArr, seeders.NewRolesSeeder(seeders.RolesSeederOptions{
			RoleRepository:       roleRepository,
			PermissionRepository: permissionRepository,
		}))
	default:
		log.Fatalf("Invalid environment: %s", env)
	}

	for _, seeder := range seedersArr {
		runSeeder(seeder.Run, seeder.Name())
	}
}
