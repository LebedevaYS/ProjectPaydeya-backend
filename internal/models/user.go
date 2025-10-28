package models

import (
    "time"
)
// User represents user model
// @Description Модель пользователя
type User struct {
    ID           int       `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    FullName     string    `json:"fullName"`
    Role         string    `json:"role"`
    AvatarURL    string    `json:"avatarUrl,omitempty"`
    IsVerified   bool      `json:"isVerified"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}
// RegisterRequest represents registration request
// @Description Запрос на регистрацию
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"fullName" binding:"required"`
    Role     string `json:"role" binding:"required,oneof=student teacher admin"`
}
// LoginRequest represents login request
// @Description Запрос на вход
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
// AuthResponse represents authentication response
// @Description Ответ с токенами и данными пользователя
type AuthResponse struct {
    Message      string `json:"message"`
    AccessToken  string `json:"accessToken,omitempty"`
    RefreshToken string `json:"refreshToken,omitempty"`
    User         *User  `json:"user,omitempty"`
}
// ForgotPasswordRequest represents forgot password request
// @Description Запрос на сброс пароля
type ForgotPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}
// ResetPasswordRequest represents reset password request
// @Description Запрос на установку нового пароля
type ResetPasswordRequest struct {
    Token       string `json:"token" binding:"required"`
    NewPassword string `json:"newPassword" binding:"required,min=6"`
}