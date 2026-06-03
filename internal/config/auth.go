package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWTSecret string
	CookieSecure bool
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 5432
	}

	return &Config{
		JWTSecret:    getEnv("JWT_SECRET", "your_super_secret_jwt_key"),
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",

		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     port,
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:     getEnv("POSTGRES_DB", "ticketing_db"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}