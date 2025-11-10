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

// GetProgress godoc
// @Summary Получить прогресс обучения
// @Description Возвращает прогресс обучения текущего пользователя
// @Tags progress
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.StudentProgress "Прогресс обучения"
// @Failure 500 {object} InternalErrorResponse "Ошибка сервера"
// @Router /student/progress [get]
func (h *ProgressHandler) GetProgress(c *gin.Context) {
    userID := c.GetInt("userID")

    progress, err := h.progressService.GetStudentProgress(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
        return
    }

    c.JSON(http.StatusOK, progress)
}

// MarkMaterialComplete godoc
// @Summary Отметить материал как завершенный
// @Description Отмечает материал как завершенный с оценкой и временем изучения
// @Tags progress
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param input body MarkCompleteRequest true "Данные завершения"
// @Success 200 {object} MarkCompleteResponse "Материал отмечен как завершенный"
// @Failure 400 {object} InvalidParametersErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} InternalErrorResponse "Ошибка сервера"
// @Router /student/materials/{id}/complete [post]
func (h *ProgressHandler) MarkMaterialComplete(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var req MarkCompleteRequest
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

// GetFavorites godoc
// @Summary Получить избранные материалы
// @Description Возвращает список избранных материалов пользователя
// @Tags progress
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} FavoritesResponse "Список избранных материалов"
// @Failure 500 {object} InternalErrorResponse "Ошибка сервера"
// @Router /student/favorites [get]
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

// ToggleFavorite godoc
// @Summary Добавить/удалить из избранного
// @Description Добавляет или удаляет материал из избранного
// @Tags progress
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param input body ToggleFavoriteRequest true "Действие с избранным"
// @Success 200 {object} ToggleFavoriteResponse "Статус избранного обновлен"
// @Failure 400 {object} InvalidParametersErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} InternalErrorResponse "Ошибка сервера"
// @Router /student/materials/{id}/favorite [post]
func (h *ProgressHandler) ToggleFavorite(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var req ToggleFavoriteRequest
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

// Request/Response models for Swagger

// MarkCompleteRequest represents mark material complete request
// @Description Запрос на отметку материала как завершенного
type MarkCompleteRequest struct {
    TimeSpent int     `json:"timeSpent" binding:"required" example:"3600"`
    Grade     float64 `json:"grade" binding:"required,min=1,max=5" example:"4.5"`
}

// MarkCompleteResponse represents mark material complete response
// @Description Ответ на отметку материала как завершенного
type MarkCompleteResponse struct {
    Message    string `json:"message" example:"Material marked as completed"`
    MaterialID int    `json:"materialID" example:"1"`
}

// ToggleFavoriteRequest represents toggle favorite request
// @Description Запрос на добавление/удаление из избранного
type ToggleFavoriteRequest struct {
    Action string `json:"action" binding:"required,oneof=add remove" example:"add"`
}

// ToggleFavoriteResponse represents toggle favorite response
// @Description Ответ на добавление/удаление из избранного
type ToggleFavoriteResponse struct {
    Message    string `json:"message" example:"Favorite updated successfully"`
    Action     string `json:"action" example:"add"`
    MaterialID int    `json:"materialID" example:"1"`
}

// FavoritesResponse represents favorites list response
// @Description Ответ со списком избранных материалов
type FavoritesResponse struct {
    Materials []interface{} `json:"materials"` // Замените на конкретный тип когда будет реализовано
    Total     int           `json:"total" example:"5"`
}

