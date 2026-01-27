package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
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
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type authService struct {
	userRepo     interfaces.UsersRepository
	emailService EmailService
}

func NewAuthService(userRepo interfaces.UsersRepository, emailService EmailService) interfaces.AuthService {
	return &authService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

var strongPasswordRegex = regexp.MustCompile(
	`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]).{8,}$`,
)

func isStrongPassword(password string) bool {
	return strongPasswordRegex.MatchString(password)
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) error {

	if err := rejectMongoOperators(req.Email); err != nil {
		return err
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	if !isStrongPassword(req.Password) {
		return errors.New(
			"password must be at least 8 characters long and contain uppercase, lowercase, number, and special character",
		)
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
		Role:              "user",
		IsActive:          false,
		ActivationToken:   activationToken,
		ActivationExpires: time.Now().Add(24 * time.Hour),
		PasswordChangedAt: time.Now(),
		PasswordExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		CreatedAt:         time.Now(),
	}

	if err := s.emailService.SendActivationEmail(user.Email, activationToken); err != nil {
		fmt.Println("Failed to send activation email, but proceeding:", err)
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *authService) ConfirmEmail(ctx context.Context, token string) error {
	fmt.Println(" SERVICE: confirming token =", token)

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

	if err := rejectMongoOperators(req.Email); err != nil {
		return err
	}
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
		fmt.Println("LOGIN ERROR: password mismatch for", user.Email)
		fmt.Println("STORED:", user.Password)
		return errors.New("invalid credentials")
	}

	otpBytes := make([]byte, 6)
	_, err = rand.Read(otpBytes)
	if err != nil {
		return err
	}
	otp := ""
	for _, b := range otpBytes {
		otp += fmt.Sprintf("%d", b%10)
	}

	user.OTP = otp
	user.OTPExpires = time.Now().Add(24 * time.Hour)

	if err := s.emailService.SendOTP(user.Email, otp); err != nil {
		fmt.Println("Failed to send OTP email, but proceeding:", err)
	}

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

	return jwt.GenerateJWT(user.ID.Hex(), user.Role)
}

func (s *authService) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error {
	const minPasswordChangeInterval = 24 * time.Hour

	// 1. Validacija user ID-a
	id, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return errors.New("invalid user id")
	}

	// 2. Učitaj korisnika
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// 3. Provera stare lozinke
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.OldPassword),
	); err != nil {
		return errors.New("stara lozinka nije ispravna")
	}

	// 4. 24h pravilo – mora proći bar 24 sata od poslednje promene
	if time.Since(user.PasswordChangedAt) < minPasswordChangeInterval {
		return errors.New("password can be changed only once every 24 hours")
	}

	// 5. (Preporuka) Nova lozinka mora biti drugačija
	if bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.NewPassword),
	) == nil {
		return errors.New("new password must be different from old password")
	}

	// 6. Hash nove lozinke
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	// 7. Update polja
	user.Password = string(hashed)
	user.PasswordChangedAt = time.Now()
	user.PasswordExpiresAt = time.Now().Add(60 * 24 * time.Hour)

	// 8. Snimi u bazu
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}

	user.ResetToken = uuid.NewString()
	user.ResetTokenExpires = time.Now().Add(10 * time.Minute)

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		fmt.Println("⚠️ Failed to update user with reset token:", err)
		return
	}

	if err := s.emailService.SendPasswordResetEmail(user.Email, user.ResetToken); err != nil {
		fmt.Println("⚠️ Failed to send password reset email:", err)
	}
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

	fmt.Println("🔄 RESET PASSWORD: New hash:", string(hashed))

	user.Password = string(hashed)
	user.PasswordChangedAt = time.Now()
	user.PasswordExpiresAt = time.Now().Add(60 * 24 * time.Hour)
	user.ResetToken = ""
	user.ResetTokenExpires = time.Time{}

	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) Logout(ctx context.Context, token string) {

}

func rejectMongoOperators(s string) error {
	if strings.Contains(s, "$") {
		return errors.New("invalid characters")
	}
	return nil
}
