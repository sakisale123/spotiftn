package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"

	"spotiftn/users/interfaces"
	"spotiftn/users/jwt"
	"spotiftn/users/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func generateToken() string {
	b := make([]byte, 32) // 256-bit token
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type authService struct {
	userRepo interfaces.UsersRepository
}

func NewAuthService(userRepo interfaces.UsersRepository) interfaces.AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// ===== REGISTER =====

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	activationToken := generateToken() // random string

	user := &models.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          string(hashed),
		IsActive:          false, // Requires email activation
		ActivationToken:   activationToken,
		ActivationExpires: time.Now().Add(24 * time.Hour),
		PasswordChangedAt: time.Now(),
		PasswordExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60 dana
		CreatedAt:         time.Now(),
	}

	// 🔴 KLJUČNO: vidi token u logu (simulacija emaila)
	fmt.Println("ACTIVATION LINK:")
	// Promenjeno da vodi na Frontend (port 3000) jer je to port koji je dozvoljen na Gateway-u (CORS)
	fmt.Println("http://localhost:3000/activate?token=" + activationToken)

	return s.userRepo.CreateUser(ctx, user)
}

// ===== EMAIL CONFIRM =====
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

// ===== LOGIN STEP 1 =====

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

	// Generisanje random 6-cifrenog koda
	otpBytes := make([]byte, 6)
	_, err = rand.Read(otpBytes)
	if err != nil {
		return err
	}
	// Mapiranje bajtova u cifre 0-9
	otp := ""
	for _, b := range otpBytes {
		otp += fmt.Sprintf("%d", b%10)
	}

	user.OTP = otp
	user.OTPExpires = time.Now().Add(24 * time.Hour)

	fmt.Println("🔐 OTP GENERATED FOR:", user.Email)
	fmt.Println("🔐 OTP CODE:", user.OTP)
	fmt.Println("🔐 OTP EXPIRES:", user.OTPExpires)

	return s.userRepo.UpdateUser(ctx, user)
}

// ===== LOGIN STEP 2 =====

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

// ===== CHANGE PASSWORD =====

func (s *authService) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error {
	id, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return errors.New("invalid user id")
	}

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if time.Since(user.PasswordChangedAt) < 24*time.Hour {
		return errors.New("password can be changed only once per day")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	user.Password = string(hashed)
	user.PasswordChangedAt = time.Now()
	user.PasswordExpiresAt = time.Now().Add(60 * 24 * time.Hour)

	return s.userRepo.UpdateUser(ctx, user)
}

// ===== FORGOT PASSWORD =====

func (s *authService) ForgotPassword(ctx context.Context, email string) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}

	user.ResetToken = uuid.NewString()
	user.ResetTokenExpires = time.Now().Add(10 * time.Minute)

	_ = s.userRepo.UpdateUser(ctx, user)
}

// ===== RESET PASSWORD =====

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

// ===== LOGOUT =====

func (s *authService) Logout(ctx context.Context, token string) {
	// noop / blacklist (nije obavezno za zadatak)
}
