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

// CreateMaterial создает новый материал
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

// GetMaterial возвращает материал по ID
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

// UpdateMaterial обновляет материал
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

// GetUserMaterials возвращает материалы пользователя
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

// PublishMaterial публикует материал
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

    // Устанавливаем значения по умолчанию если не переданы
    if req.Visibility == "" {
        req.Visibility = "published"
    }
    if req.Access == "" {
        req.Access = "open"
    }

    // TODO: реализовать в сервисе настоящую публикацию
    // Пока заглушка с генерацией shareUrl
    shareUrl := "/material/material-" + strconv.Itoa(materialID)
    if req.Access == "link" {
        // Генерируем уникальный URL для доступа по ссылке
        shareUrl = "/m/" + generateUniqueHash()
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Material published successfully",
        "materialID": materialID,
        "userID": userID,
        "visibility": req.Visibility,
        "access": req.Access,
        "shareUrl": shareUrl,
    })
}

// Вспомогательная функция для генерации хеша
func generateUniqueHash() string {
    bytes := make([]byte, 8)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}