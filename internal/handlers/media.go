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
        "message":  "Image uploaded successfully",
        "url":      result.URL,
        "fileName": result.FileName,
        "size":     result.Size,
        "width":    result.Width,
        "height":   result.Height,
    })
}

// UploadVideo загружает видео файл
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
        "message":  "Video uploaded successfully",
        "url":      result.URL,
        "fileName": result.FileName,
        "size":     result.Size,
    })
}

// EmbedVideo обрабатывает вставку видео по ссылке (YouTube, VK и т.д.)
func (h *MediaHandler) EmbedVideo(c *gin.Context) {
    var request struct {
        URL string `json:"url" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Пока простой парсинг URL
    embedURL, err := parseVideoURL(request.URL)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported video service or invalid URL"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":   "Video embed URL generated",
        "originalUrl": request.URL,
        "embedUrl":  embedURL,
        "type":      getVideoServiceType(request.URL),
    })
}

// parseVideoURL преобразует URL видео в embed URL
func parseVideoURL(url string) (string, error) {
    // TODO: Реализовать парсинг для разных видеосервисов
    // YouTube, VK Video, Vimeo и т.д.

    // Временная заглушка - возвращаем оригинальный URL
    return url, nil
}

// getVideoServiceType определяет тип видеосервиса
func getVideoServiceType(url string) string {
    // TODO: Определить сервис по домену
    return "unknown"
}

