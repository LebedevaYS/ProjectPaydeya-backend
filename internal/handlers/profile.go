package handlers

import (
    "net/http"



    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type ProfileHandler struct {
    authService *services.AuthService
}

func NewProfileHandler(authService *services.AuthService) *ProfileHandler {
    return &ProfileHandler{authService: authService}
}

// GetProfile возвращает данные профиля
func (h *ProfileHandler) GetProfile(c *gin.Context) {
    userID := c.GetInt("userID")
    userEmail := c.GetString("userEmail")
    userRole := c.GetString("userRole")

    // Пока возвращаем заглушку с данными из токена
    // Позже добавим запрос к БД
    c.JSON(http.StatusOK, gin.H{
        "id":       userID,
        "email":    userEmail,
        "role":     userRole,
        "fullName": "Полное имя из БД", // TODO: получить из БД
        "avatarUrl": "",
        "specializations": []string{},
        "message": "Это защищенный эндпоинт профиля",
    })
}

// UpdateProfile обновляет данные профиля
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
    userID := c.GetInt("userID")

    var req struct {
        FullName       string   `json:"fullName"`
        Specializations []string `json:"specializations"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Profile updated successfully",
        "userID":  userID,
        "data":    req,
    })
}

// UploadAvatar загружает аватар
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
    userID := c.GetInt("userID")

    // Пока заглушка для загрузки файла
    c.JSON(http.StatusOK, gin.H{
        "message": "Avatar upload endpoint",
        "userID":  userID,
        "avatarUrl": "/avatars/default.jpg", // TODO: реализовать загрузку
    })
}