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

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя и возвращает токены
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} models.AuthResponse "Пользователь создан"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /auth/register [post]
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

// Login godoc
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает токены
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.AuthResponse "Успешный вход"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 401 {object} ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /auth/login [post]
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

// Refresh godoc
// @Summary Обновление токенов
// @Description Обновляет access и refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RefreshTokenRequest true "Refresh токен"
// @Success 200 {object} models.AuthResponse "Токены обновлены"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 401 {object} ErrorResponse "Невалидный токен"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
    var req RefreshTokenRequest

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

// Logout godoc
// @Summary Выход из системы
// @Description Завершает сеанс пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} SuccessResponse "Успешный выход"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
    // Пока просто возвращаем успех - в будущем можно добавить blacklist токенов
    c.JSON(http.StatusOK, gin.H{
        "message": "Successfully logged out",
    })
}

// ForgotPassword godoc
// @Summary Запрос сброса пароля
// @Description Отправляет инструкции по сбросу пароля на email
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.ForgotPasswordRequest true "Email для сброса пароля"
// @Success 200 {object} SuccessResponse "Инструкции отправлены"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /auth/forgot-password [post]
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

// ResetPassword godoc
// @Summary Сброс пароля
// @Description Устанавливает новый пароль по токену сброса
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.ResetPasswordRequest true "Данные для сброса пароля"
// @Success 200 {object} SuccessResponse "Пароль изменен"
// @Failure 400 {object} ErrorResponse "Неверные данные или токен"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /auth/reset-password [post]
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

// RefreshTokenRequest represents refresh token request
// @Description Запрос на обновление токенов
type RefreshTokenRequest struct {
    RefreshToken string `json:"refreshToken" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
