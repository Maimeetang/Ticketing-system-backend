package config

import (
	"os"
	"strconv"
)

type Config struct {
	JWTSecret 	 string
	CookieSecure bool
	DBHost     	 string
	DBPort     	 int
	DBUser     	 string
	DBPassword 	 string
	DBName     	 string
}

func LoadConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 8080
	}

	return &Config{
		JWTSecret:    os.Getenv("JWT_SECRET"),
		CookieSecure: os.Getenv("APP_ENV") != "development ",
		DBHost:       os.Getenv("POSTGRES_HOST"),
		DBPort:       port,
		DBUser:       os.Getenv("POSTGRES_USER"),
		DBPassword:   os.Getenv("POSTGRES_PASSWORD"),
		DBName:       os.Getenv("POSTGRES_DB"),
	}
}