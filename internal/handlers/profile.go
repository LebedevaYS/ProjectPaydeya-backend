package handlers

import (
    "net/http"
    "strconv"
    "log"

    "paydeya-backend/internal/repositories"
    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type ProfileHandler struct {
    authService *services.AuthService
    userRepo    *repositories.UserRepository
}

func NewProfileHandler(authService *services.AuthService, userRepo *repositories.UserRepository) *ProfileHandler {
    return &ProfileHandler{
        authService: authService,
        userRepo:    userRepo,
    }
}

// GetProfile возвращает данные профиля
func (h *ProfileHandler) GetProfile(c *gin.Context) {
    userID := c.GetInt("userID")

    // Получаем пользователя из БД
    user, err := h.userRepo.GetUserProfile(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
        return
    }

    if user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Получаем специализации из БД
    specializations, err := h.userRepo.GetUserSpecializations(c.Request.Context(), userID)
    if err != nil {
        // Логируем ошибку но продолжаем (специализации не критичны)
        log.Printf("Warning: failed to get specializations for user %d: %v", userID, err)
        specializations = []string{}
    }

    c.JSON(http.StatusOK, gin.H{
        "id":               user.ID,
        "email":            user.Email,
        "fullName":         user.FullName,
        "role":             user.Role,
        "avatarUrl":        user.AvatarURL,
        "isVerified":       user.IsVerified,
        "specializations":  specializations,
        "createdAt":        user.CreatedAt,
        "updatedAt":        user.UpdatedAt,
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

    // Обновляем данные в БД
    err := h.userRepo.UpdateUserProfile(c.Request.Context(), userID, req.FullName, req.Specializations)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
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

    // Пока заглушка - в реальности сохраняем файл и получаем URL
    avatarURL := "/avatars/user_" + strconv.Itoa(userID) + ".jpg"

    // Обновляем аватар в БД
    err := h.userRepo.UpdateUserAvatar(c.Request.Context(), userID, avatarURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":   "Avatar uploaded successfully",
        "avatarUrl": avatarURL,
    })
}