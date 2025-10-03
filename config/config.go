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

func loadConfig() *Config {
	cfg := &Config{}

	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf("failed to load environment file: %w", err))
	}

	cfg.Env = os.Getenv("ENV")
	cfg.ApiVersion = os.Getenv("API_VERSION")
	cfg.AppPort = os.Getenv("APP_PORT")
	cfg.PostgresHost = os.Getenv("POSTGRES_HOST")
	cfg.PostgresPort = os.Getenv("POSTGRES_PORT")
	cfg.PostgresUsername = os.Getenv("POSTGRES_USER")
	cfg.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	cfg.PostgresName = os.Getenv("POSTGRES_DATABASE")
	cfg.MailHost = os.Getenv("MAIL_HOST")
	cfg.MailPort = os.Getenv("MAIL_PORT")
	cfg.MailUsername = os.Getenv("MAIL_USERNAME")
	cfg.MailPassword = os.Getenv("MAIL_PASSWORD")
	cfg.MailFrom = os.Getenv("MAIL_FROM")
	cfg.JwtSecret = os.Getenv("JWT_SECRET")
	cfg.StorageLocalBasePath = os.Getenv("STORAGE_LOCAL_BASE_PATH")
	expirationMinutes, err := strconv.ParseInt(os.Getenv("JWT_EXPIRATION_MINUTES"), 10, 64)
	if err != nil {
		panic(fmt.Errorf("failed to parse JWT_EXPIRATION_TIME: %w", err))
	}
	cfg.JwtExpirationMinutes = expirationMinutes
	cfg.StorageLocalBasePath = os.Getenv("STORAGE_LOCAL_BASE_PATH")

	return cfg
}

var lock = &sync.Mutex{}

var instance *Config = loadConfig()

func GetConfig() *Config {
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
