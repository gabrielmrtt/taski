package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/gabrielmrtt/taski/config"
)

func up(env string) {
	switch env {
	case "default":
		connectionUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.GetConfig().PostgresUsername,
			config.GetConfig().PostgresPassword,
			config.GetConfig().PostgresHost,
			config.GetConfig().PostgresPort,
			config.GetConfig().PostgresName,
		)
		cmd := exec.Command("migrate", "-path", "internal/core/database/migrations", "-database", connectionUrl, "up")
		stdout, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Error running migrate up: %v", err)
		}

		log.Println(string(stdout))
	case "test":
		return
	default:
		log.Fatalf("Invalid environment: %s", env)
	}
}

func down(env string, step int) {
	switch env {
	case "default":
		connectionUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.GetConfig().PostgresUsername,
			config.GetConfig().PostgresPassword,
			config.GetConfig().PostgresHost,
			config.GetConfig().PostgresPort,
			config.GetConfig().PostgresName,
		)
		cmd := exec.Command("migrate", "-path", "internal/core/database/migrations", "-database", connectionUrl, "down", fmt.Sprintf("%d", step))
		stdout, err := cmd.Output()
		if err != nil {
			log.Fatalf("Error running migrate down: %v", err)
		}

		log.Println(string(stdout))
	case "test":
		return
	default:
		log.Fatalf("Invalid environment: %s", env)
	}
}

func main() {
	var env string = "default"
	var downStep int = 0
	var method string = "up"

	if len(os.Args) >= 2 {
		env = os.Args[1]

		if len(os.Args) >= 3 {
			method = os.Args[2]

			if method == "down" && len(os.Args) < 4 {
				log.Fatalf("Step is required for down method")
			} else {
				if len(os.Args) >= 4 {
					downStepInt, err := strconv.Atoi(os.Args[3])
					if err != nil {
						log.Fatalf("Invalid step: %v", err)
					}
					downStep = downStepInt
				}
			}
		}
	}

	if method == "up" {
		up(env)
	} else {
		down(env, downStep)
	}
}
