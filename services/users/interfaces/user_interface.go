package interfaces

import (
	"context"

	"spotiftn/users/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByResetToken(ctx context.Context, token string) (*models.User, error)
	GetUserByActivationToken(ctx context.Context, token string) (*models.User, error)
}
