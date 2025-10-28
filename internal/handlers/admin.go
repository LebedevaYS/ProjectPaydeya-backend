package handlers

import (
    "net/http"
    "strconv"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type AdminHandler struct {
    adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
    return &AdminHandler{adminService: adminService}
}

// GetStatistics godoc
// @Summary Получить статистику платформы
// @Description Возвращает общую статистику платформы (пользователи, материалы, преподаватели)
// @Tags admin
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.AdminStats "Статистика платформы"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /admin/statistics [get]
func (h *AdminHandler) GetStatistics(c *gin.Context) {
    stats, err := h.adminService.GetPlatformStats(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
        return
    }

    c.JSON(http.StatusOK, stats)
}

// GetUsers godoc
// @Summary Получить список пользователей
// @Description Возвращает список пользователей с пагинацией и фильтрацией по роли
// @Tags admin
// @Produce json
// @Security ApiKeyAuth
// @Param role query string false "Фильтр по роли" Enums(student, teacher, admin)
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(20)
// @Success 200 {object} UsersListResponse "Список пользователей"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
    role := c.Query("role")
    page, _ := strconv.Atoi(c.Query("page"))
    limit, _ := strconv.Atoi(c.Query("limit"))

    if page == 0 {
        page = 1
    }
    if limit == 0 {
        limit = 20
    }

    users, total, err := h.adminService.GetUsers(c.Request.Context(), role, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "users": users,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// BlockUser godoc
// @Summary Заблокировать пользователя
// @Description Блокирует пользователя по ID с указанием причины
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID пользователя"
// @Param input body models.BlockUserRequest true "Данные для блокировки"
// @Success 200 {object} SuccessResponse "Пользователь заблокирован"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /admin/users/{id}/block [post]
func (h *AdminHandler) BlockUser(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var req models.BlockUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.adminService.BlockUser(c.Request.Context(), userID, req.Reason)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "User blocked successfully",
        "userId":  userID,
    })
}

// CreateSubject godoc
// @Summary Создать предмет
// @Description Создает новый учебный предмет
// @Tags admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body models.CreateSubjectRequest true "Данные предмета"
// @Success 200 {object} SuccessResponse "Предмет создан"
// @Failure 400 {object} ErrorResponse "Неверные данные"
// @Failure 401 {object} ErrorResponse "Требуется авторизация"
// @Failure 403 {object} ErrorResponse "Доступ запрещен"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /admin/subjects [post]
func (h *AdminHandler) CreateSubject(c *gin.Context) {
    var req models.CreateSubjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.adminService.CreateSubject(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subject"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Subject created successfully",
        "subject": req,
    })
}
// Response models for Swagger

// ErrorResponse represents error response
// @Description Стандартный ответ с ошибкой
type ErrorResponse struct {
    Error string `json:"error" example:"Authorization header required"`
}

// SuccessResponse represents success response
// @Description Стандартный успешный ответ
type SuccessResponse struct {
    Message string `json:"message" example:"Operation completed successfully"`
}

// UsersListResponse represents users list response
// @Description Ответ со списком пользователей
type UsersListResponse struct {
    Users []models.UserManagement `json:"users"`
    Total int                     `json:"total" example:"150"`
    Page  int                     `json:"page" example:"1"`
    Limit int                     `json:"limit" example:"20"`
}