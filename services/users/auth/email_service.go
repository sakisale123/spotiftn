package auth

import (
	"fmt"
	"net/smtp"
)

type EmailService interface {
	SendOTP(to string, otp string) error
	SendActivationEmail(to string, token string) error
	SendPasswordResetEmail(to string, token string) error
}

type smtpEmailService struct {
	host     string
	port     string
	email    string
	password string
}

func NewEmailService(host, port, email, password string) EmailService {
	return &smtpEmailService{
		host:     host,
		port:     port,
		email:    email,
		password: password,
	}
}

func (s *smtpEmailService) SendOTP(to string, otp string) error {
	addr := s.host + ":" + s.port
	auth := smtp.PlainAuth("", s.email, s.password, s.host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Spotiftn OTP Code\r\n"+
		"\r\n"+
		"Your OTP code is: %s\r\n", to, otp))

	if err := smtp.SendMail(addr, auth, s.email, []string{to}, msg); err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", to, err)
		return err
	}

	fmt.Printf("Email sent to %s\n", to)
	return nil
}

func (s *smtpEmailService) SendActivationEmail(to string, token string) error {
	addr := s.host + ":" + s.port
	auth := smtp.PlainAuth("", s.email, s.password, s.host)

	link := "http://localhost:3000/activate?token=" + token

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Activate Your Spotiftn Account\r\n"+
		"\r\n"+
		"Welcome to Spotiftn!\r\n"+
		"Please click the link below to activate your account:\r\n"+
		"%s\r\n", to, link))

	if err := smtp.SendMail(addr, auth, s.email, []string{to}, msg); err != nil {
		fmt.Printf(" Failed to send activation email to %s: %v\n", to, err)
		return err
	}

	fmt.Printf("Activation email sent to %s\n", to)
	return nil
}

func (s *smtpEmailService) SendPasswordResetEmail(to string, token string) error {
	addr := s.host + ":" + s.port
	auth := smtp.PlainAuth("", s.email, s.password, s.host)

	link := "http://localhost:3000/reset-password?token=" + token

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Password Reset Request\r\n"+
		"\r\n"+
		"You requested a password reset for your Spotiftn account.\r\n"+
		"Click the link below to verify your email and set a new password:\r\n"+
		"%s\r\n"+
		"\r\n"+
		"This link will expire in 10 minutes.\r\n"+
		"If you did not request this, please ignore this email.\r\n", to, link))

	if err := smtp.SendMail(addr, auth, s.email, []string{to}, msg); err != nil {
		fmt.Printf(" Failed to send reset email to %s: %v\n", to, err)
		return err
	}

	fmt.Printf(" Reset email sent to %s\n", to)
	return nil
}
