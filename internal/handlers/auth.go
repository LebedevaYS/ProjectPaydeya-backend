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