package handlers

import (
    "net/http"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/services"

    "github.com/gin-gonic/gin"
)

type CatalogHandler struct {
    catalogService *services.CatalogService
}

func NewCatalogHandler(catalogService *services.CatalogService) *CatalogHandler {
    return &CatalogHandler{catalogService: catalogService}
}

// SearchMaterials поиск материалов в каталоге
func (h *CatalogHandler) SearchMaterials(c *gin.Context) {
    var filters models.CatalogFilters

    // Парсим query параметры
    if err := c.ShouldBindQuery(&filters); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Устанавливаем значения по умолчанию для пагинации
    if filters.Page == 0 {
        filters.Page = 1
    }
    if filters.Limit == 0 {
        filters.Limit = 20
    }

    materials, total, err := h.catalogService.SearchMaterials(c.Request.Context(), filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search materials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "materials": materials,
        "total":     total,
        "page":      filters.Page,
        "limit":     filters.Limit,
        "hasMore":   (filters.Page * filters.Limit) < total,
    })
}

// GetSubjects возвращает список предметов
func (h *CatalogHandler) GetSubjects(c *gin.Context) {
    subjects, err := h.catalogService.GetSubjects(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subjects"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "subjects": subjects,
    })
}

// SearchTeachers поиск преподавателей
func (h *CatalogHandler) SearchTeachers(c *gin.Context) {
    var filters models.TeacherFilters

    if err := c.ShouldBindQuery(&filters); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    teachers, err := h.catalogService.SearchTeachers(c.Request.Context(), filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search teachers"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "teachers": teachers,
    })
}