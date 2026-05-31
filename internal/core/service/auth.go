package service

import (
	"ticketing-system/internal/apperror"
	"ticketing-system/internal/config"
	"ticketing-system/internal/core/port"
	"ticketing-system/internal/core/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil || user == nil  {
		return "", apperror.NewUnauthorized("username หรือ password ไม่ถูกต้อง")
	}

	// Check password
	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", apperror.NewUnauthorized("username หรือ password ไม่ถูกต้อง")
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
		return "", apperror.NewInternalServerError("ไม่สามารถสร้างโทเค็นได้")
	}

	return tokenString, nil
}
