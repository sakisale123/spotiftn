package auth

import (
	"context"
	"errors"
	"spotiftn/users/interfaces"
	"spotiftn/users/jwt"
	"spotiftn/users/models"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo interfaces.UsersRepository
}

func NewAuthService(userRepo interfaces.UsersRepository) interfaces.AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) error {
	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := jwt.GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}
