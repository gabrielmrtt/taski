package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/gabrielmrtt/taski/config"
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
)

type MigrationConfig struct {
	Environment string
	Method      string
	Step        int
}

type DatabaseConfig struct {
	ConnectionURL string
	Type          string
}

func getDatabaseConfig(env string) DatabaseConfig {
	switch env {
	case "default":
		coredatabase.GetPostgresConnection()
		return DatabaseConfig{
			ConnectionURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				config.GetConfig().PostgresUsername,
				config.GetConfig().PostgresPassword,
				config.GetConfig().PostgresHost,
				config.GetConfig().PostgresPort,
				config.GetConfig().PostgresName,
			),
			Type: "postgres",
		}
	case "test":
		coredatabase.GetSQLiteConnection()
		return DatabaseConfig{
			ConnectionURL: fmt.Sprintf("sqlite3://%s", "test.db"),
			Type:          "sqlite",
		}
	default:
		log.Fatalf("Invalid environment: %s", env)
		return DatabaseConfig{}
	}
}

func executeMigration(config MigrationConfig) {
	dbConfig := getDatabaseConfig(config.Environment)

	log.Printf("Executing migration %s for environment: %s", config.Method, config.Environment)
	log.Printf("Database type: %s", dbConfig.Type)
	log.Printf("Connection URL: %s", dbConfig.ConnectionURL)

	var cmd *exec.Cmd

	if config.Method == "up" {
		cmd = exec.Command("migrate", "-path", "internal/core/database/migrations", "-database", dbConfig.ConnectionURL, "up")
	} else {
		cmd = exec.Command("migrate", "-path", "internal/core/database/migrations", "-database", dbConfig.ConnectionURL, "down", strconv.Itoa(config.Step))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Running command: %s", cmd.String())

	if err := cmd.Run(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("Migration %s completed successfully", config.Method)
}

func parseArguments() MigrationConfig {
	config := MigrationConfig{
		Environment: "default",
		Method:      "up",
		Step:        0,
	}

	if len(os.Args) >= 2 {
		config.Environment = os.Args[1]
	}

	if len(os.Args) >= 3 {
		config.Method = os.Args[2]

		if config.Method == "down" && len(os.Args) < 4 {
			log.Fatalf("Step is required for down method")
		}

		if len(os.Args) >= 4 {
			step, err := strconv.Atoi(os.Args[3])
			if err != nil {
				log.Fatalf("Invalid step: %v", err)
			}
			config.Step = step
		}
	}

	return config
}

func main() {
	config := parseArguments()
	executeMigration(config)
}
