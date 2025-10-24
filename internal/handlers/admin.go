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

// GetStatistics возвращает статистику платформы
func (h *AdminHandler) GetStatistics(c *gin.Context) {
    stats, err := h.adminService.GetPlatformStats(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
        return
    }

    c.JSON(http.StatusOK, stats)
}

// GetUsers возвращает список пользователей
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

// BlockUser блокирует пользователя
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

// CreateSubject создает новый предмет
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