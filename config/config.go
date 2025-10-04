package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                  string
	ApiVersion           string
	AppPort              string
	PostgresHost         string
	PostgresPort         string
	PostgresUsername     string
	PostgresPassword     string
	PostgresName         string
	MailHost             string
	MailPort             string
	MailUsername         string
	MailPassword         string
	MailFrom             string
	JwtSecret            string
	JwtExpirationMinutes int64
	StorageLocalBasePath string
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func loadConfig() *Config {
	cfg := &Config{}

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
		fmt.Println("Continuing with system environment variables...")
	}

	cfg.Env = getEnvOrDefault("ENV", "development")
	cfg.ApiVersion = getEnvOrDefault("API_VERSION", "v1")
	cfg.AppPort = getEnvOrDefault("APP_PORT", "8080")
	cfg.PostgresHost = getEnvOrDefault("POSTGRES_HOST", "localhost")
	cfg.PostgresPort = getEnvOrDefault("POSTGRES_PORT", "5432")
	cfg.PostgresUsername = getEnvOrDefault("POSTGRES_USER", "postgres")
	cfg.PostgresPassword = getEnvOrDefault("POSTGRES_PASSWORD", "password")
	cfg.PostgresName = getEnvOrDefault("POSTGRES_DATABASE", "taski")
	cfg.MailHost = getEnvOrDefault("MAIL_HOST", "smtp.gmail.com")
	cfg.MailPort = getEnvOrDefault("MAIL_PORT", "587")
	cfg.MailUsername = getEnvOrDefault("MAIL_USERNAME", "")
	cfg.MailPassword = getEnvOrDefault("MAIL_PASSWORD", "")
	cfg.MailFrom = getEnvOrDefault("MAIL_FROM", "noreply@taski.com")
	cfg.JwtSecret = getEnvOrDefault("JWT_SECRET", "default-jwt-secret-for-development")
	cfg.StorageLocalBasePath = getEnvOrDefault("STORAGE_LOCAL_BASE_PATH", "./storage")

	expirationMinutesStr := getEnvOrDefault("JWT_EXPIRATION_MINUTES", "60")
	expirationMinutes, err := strconv.ParseInt(expirationMinutesStr, 10, 64)
	if err != nil {
		panic(fmt.Errorf("failed to parse JWT_EXPIRATION_MINUTES: %w", err))
	}
	cfg.JwtExpirationMinutes = expirationMinutes

	return cfg
}

var lock = &sync.Mutex{}

var instance *Config = loadConfig()

func GetInstance() *Config {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()

		if instance == nil {
			instance = loadConfig()
			return instance
		}
	}

	return instance
}
