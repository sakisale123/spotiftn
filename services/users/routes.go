package main

import (
	"net/http"

	"spotiftn/users/handlers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	authHandler *handlers.AuthHandler,
) {
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/confirm", authHandler.ConfirmEmail)

	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/verify-otp", authHandler.VerifyOTP)

	mux.HandleFunc("/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("/auth/change-password", authHandler.ChangePassword)

	mux.HandleFunc("/auth/logout", authHandler.Logout)
}
