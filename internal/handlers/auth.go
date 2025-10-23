package handlers

import (

    "net/http"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// Register обрабатывает регистрацию
func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.authService.Register(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Генерируем токены
    accessToken, refreshToken, err := h.authService.GenerateTokens(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
        return
    }

    c.JSON(http.StatusCreated, models.AuthResponse{
        Message:      "User created successfully",
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         user,
    })
}

// Login обрабатывает вход
func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.authService.Login(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Генерируем токены
    accessToken, refreshToken, err := h.authService.GenerateTokens(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
        return
    }

    c.JSON(http.StatusOK, models.AuthResponse{
        Message:      "Login successful",
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         user,
    })
}
// Refresh обновляет токены ← ДОБАВЬ ЭТУ ФУНКЦИЮ В handlers
func (h *AuthHandler) Refresh(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refreshToken" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    accessToken, refreshToken, err := h.authService.RefreshTokens(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, models.AuthResponse{
        Message:      "Tokens refreshed successfully",
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    })
}

// Logout выполняет выход
func (h *AuthHandler) Logout(c *gin.Context) {
    // Пока просто возвращаем успех - в будущем можно добавить blacklist токенов
    c.JSON(http.StatusOK, gin.H{
        "message": "Successfully logged out",
    })
}

// ForgotPassword обрабатывает запрос сброса пароля
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
    var req models.ForgotPasswordRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.authService.ForgotPassword(c.Request.Context(), req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Instructions sent to email",
    })
}

// ResetPassword обрабатывает сброс пароля
func (h *AuthHandler) ResetPassword(c *gin.Context) {
    var req models.ResetPasswordRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Password successfully changed",
    })
}