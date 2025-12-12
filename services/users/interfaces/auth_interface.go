package interfaces

import (
	"context"
	"spotiftn/users/models"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) error
	Login(ctx context.Context, req *models.LoginRequest) (string, error)
}
