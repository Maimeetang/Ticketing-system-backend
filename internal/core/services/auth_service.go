package services

import (
	"context"
	e "ticketing-system/internal/core/error"
	r "ticketing-system/internal/core/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(
		ctx context.Context, 
		username, 
		password string,
		JWTSecret string,
	) (string, error)
}

type authServiceImpl struct {
	userRepo r.UserRepository
}

func NewAuthService(userRepo r.UserRepository) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
	}
}

func (s *authServiceImpl) Login(
	ctx context.Context, 
	username string, 
	password string,
	JWTSecret string,
) (string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil || user == nil  {
		return "", e.NewUnauthorized("invalid username or password")
	}

	err = comparePassword(password, user.Password)
	if err != nil {
		return "", e.NewUnauthorized("invalid username or password")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", e.NewInternalServerError("unable token creation")
	}

	return tokenString, nil
}
