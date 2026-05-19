package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/config"
	"ticketing-system/internal/core/port"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo port.UserRepository
	cfg *config.AuthConfig
}

func NewAuthService(userRepo port.UserRepository, cfg *config.AuthConfig) port.AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		cfg: cfg,
	}
}

func (s *authServiceImpl) Login(username, password string) (string, error) {
	// Get user from repository
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", apperror.NewUnauthorized("The username or password is incorrect.")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", apperror.NewUnauthorized("The username or password is incorrect.")
	}

	// Create JWT Token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", apperror.NewInternalServerError("Unable to generate a token.")
	}

	return tokenString, nil
}
