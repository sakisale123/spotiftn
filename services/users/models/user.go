package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`

	Password string `bson:"password" json:"-"`

	IsActive bool `bson:"isActive"`

	CreatedAt time.Time `bson:"createdAt"`

	// Password security
	PasswordChangedAt time.Time `bson:"passwordChangedAt"`
	PasswordExpiresAt time.Time `bson:"passwordExpiresAt"`

	// Email confirmation
	ActivationToken string    `bson:"activationToken,omitempty"`
	ActivationExp   time.Time `bson:"activationExp,omitempty"`

	// OTP
	OTP        string    `bson:"otp,omitempty"`
	OTPExpires time.Time `bson:"otpExpires,omitempty"`

	// Reset password
	ResetToken        string    `bson:"resetToken,omitempty"`
	ResetTokenExpires time.Time `bson:"resetTokenExpires,omitempty"`
}
