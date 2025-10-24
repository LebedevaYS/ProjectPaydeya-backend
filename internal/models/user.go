package models

import (
    "time"
)

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

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    FullName string `json:"fullName" binding:"required"`
    Role     string `json:"role" binding:"required,oneof=student teacher admin"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Message      string `json:"message"`
    AccessToken  string `json:"accessToken,omitempty"`
    RefreshToken string `json:"refreshToken,omitempty"`
    User         *User  `json:"user,omitempty"`
}

type ForgotPasswordRequest struct {
    Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
    Token       string `json:"token" binding:"required"`
    NewPassword string `json:"newPassword" binding:"required,min=6"`
}