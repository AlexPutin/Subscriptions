package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseHost     string
	DatabasePort     int
	DatabaseURL      string
	Environment      string // e.g., "dev", "prod"
	ServerAddress    string
}

var config *Config

func MustLoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Default().Println("Error loading .env file")
	}

	dbPort, err := strconv.Atoi(MustGetEnv("DB_PORT"))
	if err != nil {
		panic(fmt.Sprintf("DB_PORT value is not integer: %s", MustGetEnv("DB_PORT")))
	}

	config = &Config{
		DatabaseUser:     MustGetEnv("DB_USER"),
		DatabasePassword: MustGetEnv("DB_PASSWORD"),
		DatabaseName:     MustGetEnv("DB_NAME"),
		DatabaseHost:     MustGetEnv("DB_HOSTNAME"),
		DatabasePort:     dbPort,
		DatabaseURL: fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			MustGetEnv("DB_USER"),
			MustGetEnv("DB_PASSWORD"),
			MustGetEnv("DB_HOSTNAME"),
			MustGetEnv("DB_PORT"),
			MustGetEnv("DB_NAME"),
		),
		Environment:   MustGetEnv("ENVIRONMENT"),
		ServerAddress: MustGetEnv("SERVER_ADDRESS"),
	}
}

func Get() *Config {
	if config == nil {
		panic("Configuration not loaded. Call MustLoadConfig() first.")
	}

	return config
}

func MustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s is not set or empty", key))
}
