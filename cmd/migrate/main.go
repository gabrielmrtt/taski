package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/gabrielmrtt/taski/config"
	shareddatabase "github.com/gabrielmrtt/taski/internal/shared/database"
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
		shareddatabase.GetPostgresConnection()
		return DatabaseConfig{
			ConnectionURL: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				config.GetInstance().PostgresUsername,
				config.GetInstance().PostgresPassword,
				config.GetInstance().PostgresHost,
				config.GetInstance().PostgresPort,
				config.GetInstance().PostgresName,
			),
			Type: "postgres",
		}
	case "test":
		return DatabaseConfig{}
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
		cmd = exec.Command(
			"migrate",
			"-path", "internal/shared/database/migrations",
			"-database", dbConfig.ConnectionURL,
			"up",
		)
	} else {
		cmd = exec.Command(
			"migrate",
			"-path", "internal/shared/database/migrations",
			"-database", dbConfig.ConnectionURL,
			"down", strconv.Itoa(config.Step),
		)
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
			log.Fatalf("Downgrade step is required when downgrading migrations")
		}

		if len(os.Args) >= 4 {
			step, err := strconv.Atoi(os.Args[3])
			if err != nil {
				log.Fatalf("Invalid downgrade step: %v", err)
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
