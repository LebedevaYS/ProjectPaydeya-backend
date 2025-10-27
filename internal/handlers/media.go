package handlers

import (
    "net/http"

    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type MediaHandler struct {
    fileService *services.FileService
}

func NewMediaHandler(fileService *services.FileService) *MediaHandler {
    return &MediaHandler{fileService: fileService}
}

// UploadImage загружает изображение
func (h *MediaHandler) UploadImage(c *gin.Context) {
    userID := c.GetInt("userID")

    file, header, err := c.Request.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
        return
    }
    defer file.Close()

    result, err := h.fileService.UploadImage(c.Request.Context(), file, header.Filename, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "url":      result.URL,
        "fileName": result.FileName,
        "size":     result.Size,
    })
}

// UploadVideo загружает видео
func (h *MediaHandler) UploadVideo(c *gin.Context) {
    userID := c.GetInt("userID")

    file, header, err := c.Request.FormFile("video")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
        return
    }
    defer file.Close()

    result, err := h.fileService.UploadVideo(c.Request.Context(), file, header.Filename, userID, header.Size)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "url":      result.URL,
        "fileName": result.FileName,
        "size":     result.Size,
    })
}

// EmbedVideo обрабатывает вставку видео по ссылке
func (h *MediaHandler) EmbedVideo(c *gin.Context) {
    var request struct {
        URL string `json:"url" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Пока заглушка - можно интегрировать с YouTube/VK Video API
    c.JSON(http.StatusOK, gin.H{
        "embedUrl": request.URL,
        "message":  "Video embedding will be implemented soon",
    })
}