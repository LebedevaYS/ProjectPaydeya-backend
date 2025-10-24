package handlers

import (
    "net/http"
    "strconv"

    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type ProgressHandler struct {
    progressService *services.ProgressService
}

func NewProgressHandler(progressService *services.ProgressService) *ProgressHandler {
    return &ProgressHandler{progressService: progressService}
}

// GetProgress возвращает прогресс ученика
func (h *ProgressHandler) GetProgress(c *gin.Context) {
    userID := c.GetInt("userID")

    progress, err := h.progressService.GetStudentProgress(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
        return
    }

    c.JSON(http.StatusOK, progress)
}

// MarkMaterialComplete отмечает материал как завершенный
func (h *ProgressHandler) MarkMaterialComplete(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var req struct {
        TimeSpent int     `json:"timeSpent" binding:"required"`
        Grade     float64 `json:"grade" binding:"required,min=1,max=5"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.progressService.MarkMaterialComplete(c.Request.Context(), userID, materialID, req.TimeSpent, req.Grade)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark material as complete"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Material marked as completed",
        "materialID": materialID,
    })
}

// GetFavorites возвращает избранные материалы
func (h *ProgressHandler) GetFavorites(c *gin.Context) {
    userID := c.GetInt("userID")

    materials, err := h.progressService.GetFavoriteMaterials(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "materials": materials,
        "total":     len(materials),
    })
}

// ToggleFavorite добавляет/удаляет материал из избранного
func (h *ProgressHandler) ToggleFavorite(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var req struct {
        Action string `json:"action" binding:"required,oneof=add remove"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.progressService.ToggleFavorite(c.Request.Context(), userID, materialID, req.Action)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorites"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Favorite updated successfully",
        "action":  req.Action,
        "materialID": materialID,
    })
}