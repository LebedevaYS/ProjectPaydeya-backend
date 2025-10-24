package services

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/repositories"
)

type MaterialService struct {
    materialRepo *repositories.MaterialRepository
    blockRepo    *repositories.BlockRepository
}

func NewMaterialService(materialRepo *repositories.MaterialRepository, blockRepo *repositories.BlockRepository) *MaterialService {
    return &MaterialService{
        materialRepo: materialRepo,
        blockRepo:    blockRepo,
    }
}

// CreateMaterial создает новый материал
func (s *MaterialService) CreateMaterial(ctx context.Context, userID int, req *models.CreateMaterialRequest) (*models.Material, error) {
    material := &models.Material{
        Title:    req.Title,
        Subject:  req.Subject,
        AuthorID: userID,
        Status:   "draft",
        Access:   "open",
        Blocks:   []models.Block{},
    }

    if err := s.materialRepo.CreateMaterial(ctx, material); err != nil {
        return nil, fmt.Errorf("failed to create material: %w", err)
    }

    return material, nil
}

// GetMaterial возвращает материал с блоками
func (s *MaterialService) GetMaterial(ctx context.Context, materialID int) (*models.Material, error) {
    material, err := s.materialRepo.GetMaterial(ctx, materialID)
    if err != nil || material == nil {
        return nil, err
    }

    // Загружаем блоки
    blocks, err := s.blockRepo.GetBlocks(ctx, materialID)
    if err != nil {
        return nil, err
    }

    material.Blocks = blocks
    return material, nil
}

// UpdateMaterial обновляет материал и блоки
func (s *MaterialService) UpdateMaterial(ctx context.Context, userID int, materialID int, req *models.UpdateMaterialRequest) error {
    // Получаем текущий материал для проверки прав
    material, err := s.materialRepo.GetMaterial(ctx, materialID)
    if err != nil || material == nil {
        return fmt.Errorf("material not found")
    }

    if material.AuthorID != userID {
        return fmt.Errorf("access denied")
    }

    // Обновляем заголовок если передан
    if req.Title != "" {
        material.Title = req.Title
        if err := s.materialRepo.UpdateMaterial(ctx, material); err != nil {
            return err
        }
    }

    // Сохраняем блоки если переданы
    if req.Blocks != nil {
        if err := s.blockRepo.SaveBlocks(ctx, materialID, req.Blocks); err != nil {
            return fmt.Errorf("failed to save blocks: %w", err)
        }
    }

    return nil
}

// generateShareURL генерирует уникальный URL для доступа по ссылке
func (s *MaterialService) generateShareURL() string {
    bytes := make([]byte, 8)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}