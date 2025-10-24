package middleware

import (
    "net/http"
    "strings"

    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Получаем токен из заголовка
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Формат: Bearer <token>
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        token := parts[1]

        // Проверяем токен
        claims, err := authService.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Сохраняем данные пользователя в контекст
        c.Set("userID", claims.UserID)
        c.Set("userEmail", claims.Email)
        c.Set("userRole", claims.Role)

        c.Next()
    }
}