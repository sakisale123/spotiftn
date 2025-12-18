package repository

import (
	"context"
	"errors"
	"fmt"
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
	fmt.Println("🧠 REPO USING DB:", db.Name())

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

func (r *usersRepository) UpdateUser(
	ctx context.Context,
	user *models.User,
) error {

	fmt.Println(
		"🟢 REPO UPDATE:",
		"isActive =", user.IsActive,
		"activationToken =", user.ActivationToken,
		"activationExpires =", user.ActivationExpires,
	)

	update := bson.M{
		"$set": bson.M{
			"isActive":          user.IsActive,
			"activationToken":   user.ActivationToken,
			"activationExpires": user.ActivationExpires,
			"password":          user.Password,
			"passwordChangedAt": user.PasswordChangedAt,
			"passwordExpiresAt": user.PasswordExpiresAt,
		},
	}

	res, err := r.collection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("no user matched for update")
	}

	fmt.Println("🟢 REPO UPDATE OK, modified:", res.ModifiedCount)

	return nil
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
func (r *usersRepository) GetUserByActivationToken(
	ctx context.Context,
	token string,
) (*models.User, error) {

	fmt.Println("🔎 REPO QUERY TOKEN =", token)

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{
		"activationToken": token,
	}).Decode(&user)

	if err != nil {
		fmt.Println("❌ REPO FIND ERROR:", err)
		return nil, err
	}

	return &user, nil
}
