package config

import "os"

type AuthConfig struct {
	JWTSecret string
	CookieSecure bool
}

func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret: os.Getenv("JWT_SECRET"),
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
	}
}