package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort   string
	AppEnv    string

	PostgresHost     string
	PostgresPort     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string

	DefaultAdminUsername string
	DefaultAdminPassword string
	DefaultAdminFirstname string
	DefaultAdminLastname string
	DefaultAdminPhonenumber string

	JwtSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config := &Config{
		AppPort:          getEnv("PORT", "8080"),
		AppEnv:           getEnv("APP_ENV", "development"),
		
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),

		DefaultAdminUsername:    os.Getenv("DEFAULT_ADMIN_USERNAME"),
		DefaultAdminPassword:    os.Getenv("DEFAULT_ADMIN_PASSWORD"),
		DefaultAdminFirstname:   os.Getenv("DEFAULT_ADMIN_FIRSTNAME"),
		DefaultAdminLastname:    os.Getenv("DEFAULT_ADMIN_LASTNAME"),
		DefaultAdminPhonenumber: os.Getenv("DEFAULT_ADMIN_PHONENUMBER"),
		
		JwtSecret:        os.Getenv("JWT_SECRET"),
	}

	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

