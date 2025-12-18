package interfaces

import (
	"context"

	"spotiftn/users/models"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) error
	ConfirmEmail(ctx context.Context, token string) error

	// Login (2-step)
	LoginStep1(ctx context.Context, req *models.LoginRequest) error
	VerifyOTP(ctx context.Context, req *models.OTPVerifyRequest) (string, error)

	// Password management
	ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, email string)
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error

	Logout(ctx context.Context, token string)
}
