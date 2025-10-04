package main

import (
	"log"
	"os"

	"github.com/gabrielmrtt/taski/cmd/seed/seeders"
	roledatabase "github.com/gabrielmrtt/taski/internal/role/infra/database"
	shareddatabase "github.com/gabrielmrtt/taski/internal/shared/database"
	"github.com/uptrace/bun"
)

type Seeder interface {
	Name() string
	Run() error
}

type SeedConfig struct {
	Environment string
}

type DatabaseRepositories struct {
	PermissionRepository *roledatabase.PermissionBunRepository
	RoleRepository       *roledatabase.RoleBunRepository
}

func getDatabaseRepositories(env string) (DatabaseRepositories, *bun.DB) {
	var connection *bun.DB
	var permissionRepository *roledatabase.PermissionBunRepository
	var roleRepository *roledatabase.RoleBunRepository

	switch env {
	case "default":
		connection = shareddatabase.GetPostgresConnection()
		permissionRepository = roledatabase.NewPermissionBunRepository(connection)
		roleRepository = roledatabase.NewRoleBunRepository(connection)
	case "test":
		log.Fatalf("Test environment is not supported for seeding")
	default:
		log.Fatalf("Invalid environment: %s", env)
	}

	return DatabaseRepositories{
		PermissionRepository: permissionRepository,
		RoleRepository:       roleRepository,
	}, connection
}

func createSeeders(repos DatabaseRepositories) []Seeder {
	var seedersArr []Seeder = make([]Seeder, 0)

	seedersArr = append(seedersArr, seeders.NewPermissionSeeder(seeders.PermissionSeederOptions{
		PermissionRepository: repos.PermissionRepository,
	}))

	seedersArr = append(seedersArr, seeders.NewRolesSeeder(seeders.RolesSeederOptions{
		RoleRepository:       repos.RoleRepository,
		PermissionRepository: repos.PermissionRepository,
	}))

	return seedersArr
}

func executeSeeder(seeder Seeder) {
	log.Printf("Starting seeder: %s", seeder.Name())

	err := seeder.Run()
	if err != nil {
		log.Fatalf("Seeder '%s' failed: %v", seeder.Name(), err)
	}

	log.Printf("Seeder '%s' completed successfully", seeder.Name())
}

func runSeeders(seedersArr []Seeder) {
	log.Printf("Running %d seeders for environment", len(seedersArr))

	for i, seeder := range seedersArr {
		log.Printf("Progress: %d/%d", i+1, len(seedersArr))
		executeSeeder(seeder)
	}

	log.Println("All seeders completed successfully")
}

func parseArguments() SeedConfig {
	config := SeedConfig{
		Environment: "default",
	}

	if len(os.Args) >= 2 {
		config.Environment = os.Args[1]
	}

	return config
}

func main() {
	config := parseArguments()

	log.Printf("Starting seed process for environment: %s", config.Environment)

	repos, connection := getDatabaseRepositories(config.Environment)
	_ = connection

	seedersArr := createSeeders(repos)
	runSeeders(seedersArr)
}
