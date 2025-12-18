package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"spotiftn/users/interfaces"
	"spotiftn/users/models"
)

type usersRepository struct {
	collection *mongo.Collection
}

func NewUsersRepository(db *mongo.Database) interfaces.UsersRepository {
	return &usersRepository{
		collection: db.Collection("users"),
	}
}

//
// ===== BASIC USER =====
//

func (r *usersRepository) CreateUser(ctx context.Context, user *models.User) error {
	var existing models.User
	err := r.collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existing)
	if err == nil {
		return errors.New("email already in use")
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, user)
	return err
}

func (r *usersRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//
// ===== GENERIC UPDATE =====
//

func (r *usersRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.UpdateByID(
		ctx,
		user.ID,
		bson.M{"$set": user},
	)
	return err
}

func (r *usersRepository) GetUserByResetToken(
	ctx context.Context,
	token string,
) (*models.User, error) {

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{
		"resetToken": token,
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
