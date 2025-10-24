package services

import (
    "context"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/repositories"
)

type ProgressService struct {
    progressRepo *repositories.ProgressRepository
}

func NewProgressService(progressRepo *repositories.ProgressRepository) *ProgressService {
    return &ProgressService{progressRepo: progressRepo}
}

// GetStudentProgress возвращает прогресс ученика
func (s *ProgressService) GetStudentProgress(ctx context.Context, userID int) (*models.StudentProgress, error) {
    return s.progressRepo.GetStudentProgress(ctx, userID)
}

// MarkMaterialComplete отмечает материал как завершенный
func (s *ProgressService) MarkMaterialComplete(ctx context.Context, userID, materialID int, timeSpent int, grade float64) error {
    return s.progressRepo.MarkMaterialComplete(ctx, userID, materialID, timeSpent, grade)
}

// GetFavoriteMaterials возвращает избранные материалы
func (s *ProgressService) GetFavoriteMaterials(ctx context.Context, userID int) ([]models.CatalogMaterial, error) {
    return s.progressRepo.GetFavoriteMaterials(ctx, userID)
}

// ToggleFavorite добавляет/удаляет материал из избранного
func (s *ProgressService) ToggleFavorite(ctx context.Context, userID, materialID int, action string) error {
    return s.progressRepo.ToggleFavorite(ctx, userID, materialID, action)
}