package services

import (
    "context"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/repositories"
)

type AdminService struct {
    adminRepo *repositories.AdminRepository
}

func NewAdminService(adminRepo *repositories.AdminRepository) *AdminService {
    return &AdminService{adminRepo: adminRepo}
}

// GetPlatformStats возвращает статистику платформы
func (s *AdminService) GetPlatformStats(ctx context.Context) (*models.AdminStats, error) {
    return s.adminRepo.GetPlatformStats(ctx)
}

// GetUsers возвращает список пользователей
func (s *AdminService) GetUsers(ctx context.Context, role string, page, limit int) ([]models.UserManagement, int, error) {
    return s.adminRepo.GetUsers(ctx, role, page, limit)
}

// BlockUser блокирует пользователя
func (s *AdminService) BlockUser(ctx context.Context, userID int, reason string) error {
    return s.adminRepo.BlockUser(ctx, userID, reason)
}

// CreateSubject создает новый предмет
func (s *AdminService) CreateSubject(ctx context.Context, req *models.CreateSubjectRequest) error {
    return s.adminRepo.CreateSubject(ctx, req)
}