package services

import (
    "context"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/repositories"
)

type CatalogService struct {
    catalogRepo *repositories.CatalogRepository
}

func NewCatalogService(catalogRepo *repositories.CatalogRepository) *CatalogService {
    return &CatalogService{catalogRepo: catalogRepo}
}

// SearchMaterials поиск материалов с фильтрацией
func (s *CatalogService) SearchMaterials(ctx context.Context, filters models.CatalogFilters) ([]models.CatalogMaterial, int, error) {
    return s.catalogRepo.SearchMaterials(ctx, filters)
}

// GetSubjects возвращает список предметов
func (s *CatalogService) GetSubjects(ctx context.Context) ([]models.Subject, error) {
    return s.catalogRepo.GetSubjects(ctx)
}

// SearchTeachers поиск преподавателей
func (s *CatalogService) SearchTeachers(ctx context.Context, filters models.TeacherFilters) ([]models.Teacher, error) {
    return s.catalogRepo.SearchTeachers(ctx, filters)
}