package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	mathrand "math/rand"
	"time"

	"encoding/hex"

	"github.com/google/uuid"

	"spotiftn/users/interfaces"
	"spotiftn/users/jwt"
	"spotiftn/users/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func generateOTP() string {
	return fmt.Sprintf("%06d", mathrand.Intn(1000000))
}

type authService struct {
	userRepo interfaces.UsersRepository
}

func NewAuthService(userRepo interfaces.UsersRepository) interfaces.AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	activationToken := generateToken()

	user := &models.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          string(hashed),
		IsActive:          false,
		ActivationToken:   activationToken,
		ActivationExpires: time.Now().Add(24 * time.Hour),
		PasswordChangedAt: time.Now(),
		PasswordExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		CreatedAt:         time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	log.Println("✅ REGISTER – TOKEN SAVED IN DB:", user.ActivationToken)

	log.Println("📩 ACTIVATION LINK:")
	log.Println("http://localhost:8081/auth/confirm?token=" + activationToken)
	return nil

}

func (s *authService) ConfirmEmail(ctx context.Context, token string) error {
	fmt.Println("🧠 SERVICE: confirming token =", token)

	user, err := s.userRepo.GetUserByActivationToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired activation token")
	}

	if time.Now().After(user.ActivationExpires) {
		return errors.New("activation token expired")
	}

	user.IsActive = true
	user.ActivationToken = ""
	user.ActivationExpires = time.Time{}

	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) LoginStep1(ctx context.Context, req *models.LoginRequest) error {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("invalid credentials")
	}

	if !user.IsActive {
		return errors.New("account not activated")
	}

	if time.Now().After(user.PasswordExpiresAt) {
		return errors.New("password expired")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return errors.New("invalid credentials")
	}
	otp := generateOTP()

	user.OTP = otp
	user.OTPExpires = time.Now().Add(5 * time.Minute)

	log.Println("🔐 OTP GENERATED:", otp)

	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) VerifyOTP(ctx context.Context, req *models.OTPVerifyRequest) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid otp")
	}

	if user.OTP != req.OTP || time.Now().After(user.OTPExpires) {
		return "", errors.New("invalid or expired otp")
	}

	user.OTP = ""
	user.OTPExpires = time.Time{}
	_ = s.userRepo.UpdateUser(ctx, user)

	return jwt.GenerateJWT(user.ID.Hex())
}

func (s *authService) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error {
	// 1. Parse user ID
	id, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return errors.New("invalid user id")
	}

	// 2. Get user
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	// 3. Verify old password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.OldPassword),
	); err != nil {
		return errors.New("old password is incorrect")
	}

	// 4. Enforce 1-day rule
	if !user.PasswordChangedAt.IsZero() &&
		time.Since(user.PasswordChangedAt) < 24*time.Hour {
		return errors.New("password can be changed only once per day")
	}

	// 5. Hash new password
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	// 6. Update fields
	user.Password = string(hashed)
	user.PasswordChangedAt = time.Now()
	user.PasswordExpiresAt = time.Now().Add(60 * 24 * time.Hour)

	// 7. Save
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}

	user.ResetToken = uuid.NewString()
	user.ResetTokenExpires = time.Now().Add(10 * time.Minute)

	_ = s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	user, err := s.userRepo.GetUserByResetToken(ctx, req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(user.ResetTokenExpires) {
		return errors.New("token expired")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	user.Password = string(hashed)
	user.PasswordChangedAt = time.Now()
	user.PasswordExpiresAt = time.Now().Add(60 * 24 * time.Hour)
	user.ResetToken = ""
	user.ResetTokenExpires = time.Time{}

	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) Logout(ctx context.Context, token string) {
	// noop / blacklist (nije obavezno za zadatak)
}
