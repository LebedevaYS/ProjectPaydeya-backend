package handlers

import (
    "net/http"
    "strconv"
    "crypto/rand"
    "encoding/hex"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type MaterialHandler struct {
    materialService *services.MaterialService
}

func NewMaterialHandler(materialService *services.MaterialService) *MaterialHandler {
    return &MaterialHandler{materialService: materialService}
}

// CreateMaterial godoc
// @Summary Создать материал
// @Description Создает новый учебный материал
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body models.CreateMaterialRequest true "Данные материала"
// @Success 201 {object} CreateMaterialResponse "Материал создан"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials [post]
func (h *MaterialHandler) CreateMaterial(c *gin.Context) {
    userID := c.GetInt("userID")

    var req models.CreateMaterialRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    material, err := h.materialService.CreateMaterial(c.Request.Context(), userID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":  "Material created successfully",
        "material": material,
        "editorUrl": "/editor/" + strconv.Itoa(material.ID),
    })
}

// GetMaterial godoc
// @Summary Получить материал
// @Description Возвращает материал по ID
// @Tags materials
// @Accept json
// @Produce json
// @Param id path int true "ID материала"
// @Success 200 {object} models.Material "Материал"
// @Failure 400 {object} ErrorResponse "Неверный ID"
// @Failure 404 {object} ErrorResponse "Материал не найден"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/{id} [get]
func (h *MaterialHandler) GetMaterial(c *gin.Context) {
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    material, err := h.materialService.GetMaterial(c.Request.Context(), materialID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if material == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
        return
    }

    c.JSON(http.StatusOK, material)
}

// UpdateMaterial godoc
// @Summary Обновить материал
// @Description Обновляет материал (только для автора)
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param input body models.UpdateMaterialRequest true "Данные для обновления"
// @Success 200 {object} SuccessResponse "Материал обновлен"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/{id} [put]
func (h *MaterialHandler) UpdateMaterial(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var req models.UpdateMaterialRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.materialService.UpdateMaterial(c.Request.Context(), userID, materialID, &req)
    if err != nil {
        if err.Error() == "access denied" {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Material updated successfully",
    })
}

// GetUserMaterials godoc
// @Summary Получить материалы пользователя
// @Description Возвращает материалы текущего пользователя
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "Фильтр по статусу" Enums(draft, published, archived)
// @Success 200 {object} UserMaterialsResponse "Материалы пользователя"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/my [get]
func (h *MaterialHandler) GetUserMaterials(c *gin.Context) {
    userID := c.GetInt("userID")
    status := c.Query("status") // draft, published, archived

    // TODO: реализовать в сервисе
    c.JSON(http.StatusOK, gin.H{
        "message": "User materials endpoint",
        "userID":  userID,
        "status":  status,
        "materials": []string{}, // заглушка
    })
}

// PublishMaterial godoc
// @Summary Опубликовать материал
// @Description Публикует материал с указанными настройками видимости
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param input body models.PublishMaterialRequest true "Настройки публикации"
// @Success 200 {object} PublishMaterialResponse "Материал опубликован"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/{id}/publish [post]
func (h *MaterialHandler) PublishMaterial(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var req models.PublishMaterialRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Устанавливаем значения по умолчанию
    if req.Visibility == "" {
        req.Visibility = "published"
    }
    if req.Access == "" {
        req.Access = "open"
    }

    // Вызываем настоящую логику публикации
    material, err := h.materialService.PublishMaterial(c.Request.Context(), userID, materialID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Material published successfully",
        "material": material,
        "shareUrl": material.ShareURL,
    })
}

// AddBlock godoc
// @Summary Добавить блок
// @Description Добавляет блок к материалу
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param input body models.Block true "Данные блока"
// @Success 200 {object} AddBlockResponse "Блок добавлен"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/{id}/blocks [post]
func (h *MaterialHandler) AddBlock(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var block models.Block
    if err := c.ShouldBindJSON(&block); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.materialService.AddBlock(c.Request.Context(), userID, materialID, &block)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Block added successfully",
        "blockId": block.ID,
    })
}

// UpdateBlock godoc
// @Summary Обновить блок
// @Description Обновляет блок материала
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param blockId path string true "ID блока"
// @Param input body models.Block true "Данные блока"
// @Success 200 {object} SuccessResponse "Блок обновлен"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/{id}/blocks/{blockId} [put]
func (h *MaterialHandler) UpdateBlock(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    blockID := c.Param("blockId")

    var block models.Block
    if err := c.ShouldBindJSON(&block); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    block.ID = blockID

    err = h.materialService.UpdateBlock(c.Request.Context(), userID, materialID, blockID, &block)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Block updated successfully",
    })
}

// DeleteBlock godoc
// @Summary Удалить блок
// @Description Удаляет блок из материала
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param blockId path string true "ID блока"
// @Success 200 {object} SuccessResponse "Блок удален"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/{id}/blocks/{blockId} [delete]
func (h *MaterialHandler) DeleteBlock(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    blockID := c.Param("blockId")

    err = h.materialService.DeleteBlock(c.Request.Context(), userID, materialID, blockID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Block deleted successfully",
    })
}

// ReorderBlocks godoc
// @Summary Изменить порядок блоков
// @Description Изменяет порядок блоков в материале
// @Tags materials
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID материала"
// @Param input body ReorderBlocksRequest true "Новый порядок блоков"
// @Success 200 {object} ReorderBlocksResponse "Порядок изменен"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 500 {object} ErrorResponse "Ошибка сервера"
// @Router /materials/{id}/reorder [put]
func (h *MaterialHandler) ReorderBlocks(c *gin.Context) {
    userID := c.GetInt("userID")
    materialID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
        return
    }

    var req ReorderBlocksRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.materialService.ReorderBlocks(c.Request.Context(), userID, materialID, req.Blocks)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Blocks reordered successfully",
        "newOrder": req.Blocks,
    })
}

// Вспомогательная функция для генерации хеша
func generateUniqueHash() string {
    bytes := make([]byte, 8)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// Response models for Swagger

// CreateMaterialResponse represents create material response
// @Description Ответ на создание материала
type CreateMaterialResponse struct {
    Message   string          `json:"message" example:"Material created successfully"`
    Material  models.Material `json:"material"`
    EditorURL string          `json:"editorUrl" example:"/editor/1"`
}

// PublishMaterialResponse represents publish material response
// @Description Ответ на публикацию материала
type PublishMaterialResponse struct {
    Message   string          `json:"message" example:"Material published successfully"`
    Material  models.Material `json:"material"`
    ShareURL  string          `json:"shareUrl" example:"https://paydeya.com/share/abc123"`
}

// AddBlockResponse represents add block response
// @Description Ответ на добавление блока
type AddBlockResponse struct {
    Message string `json:"message" example:"Block added successfully"`
    BlockID string `json:"blockId" example:"block_123"`
}

// UserMaterialsResponse represents user materials response
// @Description Ответ со списком материалов пользователя
type UserMaterialsResponse struct {
    Message    string            `json:"message" example:"User materials endpoint"`
    UserID     int               `json:"userID" example:"123"`
    Status     string            `json:"status" example:"draft"`
    Materials  []string          `json:"materials"`
}

// ReorderBlocksRequest represents reorder blocks request
// @Description Запрос на изменение порядка блоков
type ReorderBlocksRequest struct {
    Blocks []string `json:"blocks" binding:"required" example:"block_1,block_2,block_3"`
}

// ReorderBlocksResponse represents reorder blocks response
// @Description Ответ на изменение порядка блоков
type ReorderBlocksResponse struct {
    Message  string   `json:"message" example:"Blocks reordered successfully"`
    NewOrder []string `json:"newOrder" example:"block_1,block_2,block_3"`
}
