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

// UploadImage godoc
// @Summary Загрузить изображение
// @Description Загружает изображение для использования в материалах
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param image formData file true "Файл изображения"
// @Success 200 {object} UploadImageResponse "Изображение загружено"
// @Failure 400 {object} ErrorResponse "Неверный файл"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки"
// @Router /media/images [post]
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

// UploadVideo godoc
// @Summary Загрузить видео
// @Description Загружает видео файл для использования в материалах
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param video formData file true "Файл видео"
// @Success 200 {object} UploadVideoResponse "Видео загружено"
// @Failure 400 {object} ErrorResponse "Неверный файл"
// @Failure 500 {object} ErrorResponse "Ошибка загрузки"
// @Router /media/videos [post]
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

// EmbedVideo godoc
// @Summary Вставить видео по ссылке
// @Description Создает embed ссылку для видео с YouTube, VK и других сервисов
// @Tags media
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body EmbedVideoRequest true "URL видео"
// @Success 200 {object} EmbedVideoResponse "Ссылка для вставки создана"
// @Failure 400 {object} ErrorResponse "Неверный URL"
// @Failure 500 {object} ErrorResponse "Ошибка обработки"
// @Router /media/embed [post]
func (h *MediaHandler) EmbedVideo(c *gin.Context) {
    var request EmbedVideoRequest

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

// Request/Response models for Swagger

// UploadImageResponse represents image upload response
// @Description Ответ на загрузку изображения
type UploadImageResponse struct {
    Message  string `json:"message" example:"Image uploaded successfully"`
    URL      string `json:"url" example:"https://example.com/images/abc123.jpg"`
    FileName string `json:"fileName" example:"image.jpg"`
    Size     int64  `json:"size" example:"102400"`
    Width    int    `json:"width" example:"1920"`
    Height   int    `json:"height" example:"1080"`
}

// UploadVideoResponse represents video upload response
// @Description Ответ на загрузку видео
type UploadVideoResponse struct {
    Message  string `json:"message" example:"Video uploaded successfully"`
    URL      string `json:"url" example:"https://example.com/videos/abc123.mp4"`
    FileName string `json:"fileName" example:"video.mp4"`
    Size     int64  `json:"size" example:"10485760"`
}

// EmbedVideoRequest represents embed video request
// @Description Запрос на создание embed ссылки для видео
type EmbedVideoRequest struct {
    URL string `json:"url" binding:"required" example:"https://www.youtube.com/watch?v=abc123"`
}

// EmbedVideoResponse represents embed video response
// @Description Ответ с embed ссылкой для видео
type EmbedVideoResponse struct {
    Message      string `json:"message" example:"Video embed URL generated"`
    OriginalURL  string `json:"originalUrl" example:"https://www.youtube.com/watch?v=abc123"`
    EmbedURL     string `json:"embedUrl" example:"https://www.youtube.com/embed/abc123"`
    Type         string `json:"type" example:"youtube"`
}

