package services

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "strconv"

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
// PublishMaterial публикует материал
func (s *MaterialService) PublishMaterial(ctx context.Context, userID, materialID int, req *models.PublishMaterialRequest) (*models.Material, error) {
    // Получаем материал
    material, err := s.materialRepo.GetMaterial(ctx, materialID)
    if err != nil || material == nil {
        return nil, fmt.Errorf("material not found")
    }

    // Проверяем права
    if material.AuthorID != userID {
        return nil, fmt.Errorf("access denied")
    }

    // Обновляем статус и доступ
    material.Status = req.Visibility
    material.Access = req.Access

    // Генерируем share URL если нужно
    if req.Access == "link" {
        material.ShareURL = "/m/" + s.generateShareURL()
    } else {
        material.ShareURL = "/material/" + strconv.Itoa(materialID)
    }

    // Сохраняем изменения
    if err := s.materialRepo.UpdateMaterial(ctx, material); err != nil {
        return nil, fmt.Errorf("failed to publish material: %w", err)
    }

    return material, nil
}

// generateShareURL генерирует уникальный URL для доступа по ссылке
func (s *MaterialService) generateShareURL() string {
    bytes := make([]byte, 8)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// AddBlock добавляет блок к материалу
func (s *MaterialService) AddBlock(ctx context.Context, userID, materialID int, block *models.Block) error {
    // Проверяем права
    material, err := s.materialRepo.GetMaterial(ctx, materialID)
    if err != nil || material == nil {
        return fmt.Errorf("material not found")
    }
    if material.AuthorID != userID {
        return fmt.Errorf("access denied")
    }

    // Получаем текущие блоки
    blocks, err := s.blockRepo.GetBlocks(ctx, materialID)
    if err != nil {
        return err
    }

    // Добавляем новый блок
    blocks = append(blocks, *block)

    // Сохраняем все блоки
    return s.blockRepo.SaveBlocks(ctx, materialID, blocks)
}

// UpdateBlock обновляет блок
func (s *MaterialService) UpdateBlock(ctx context.Context, userID, materialID int, blockID string, block *models.Block) error {
    // Проверяем права
    material, err := s.materialRepo.GetMaterial(ctx, materialID)
    if err != nil || material == nil {
        return fmt.Errorf("material not found")
    }
    if material.AuthorID != userID {
        return fmt.Errorf("access denied")
    }

    // Получаем текущие блоки
    blocks, err := s.blockRepo.GetBlocks(ctx, materialID)
    if err != nil {
        return err
    }

    // Находим и обновляем блок
    for i, b := range blocks {
        if b.ID == blockID {
            blocks[i] = *block
            return s.blockRepo.SaveBlocks(ctx, materialID, blocks)
        }
    }

    return fmt.Errorf("block not found")
}

// DeleteBlock удаляет блок
func (s *MaterialService) DeleteBlock(ctx context.Context, userID, materialID int, blockID string) error {
    // Проверяем права
    material, err := s.materialRepo.GetMaterial(ctx, materialID)
    if err != nil || material == nil {
        return fmt.Errorf("material not found")
    }
    if material.AuthorID != userID {
        return fmt.Errorf("access denied")
    }

    // Получаем текущие блоки
    blocks, err := s.blockRepo.GetBlocks(ctx, materialID)
    if err != nil {
        return err
    }

    // Удаляем блок
    var newBlocks []models.Block
    for _, block := range blocks {
        if block.ID != blockID {
            newBlocks = append(newBlocks, block)
        }
    }

    return s.blockRepo.SaveBlocks(ctx, materialID, newBlocks)
}

// ReorderBlocks изменяет порядок блоков
func (s *MaterialService) ReorderBlocks(ctx context.Context, userID, materialID int, blockIDs []string) error {
    // Проверяем права
    material, err := s.materialRepo.GetMaterial(ctx, materialID)
    if err != nil || material == nil {
        return fmt.Errorf("material not found")
    }
    if material.AuthorID != userID {
        return fmt.Errorf("access denied")
    }

    // Получаем текущие блоки
    blocks, err := s.blockRepo.GetBlocks(ctx, materialID)
    if err != nil {
        return err
    }

    // Создаем мапу для быстрого поиска блоков по ID
    blockMap := make(map[string]models.Block)
    for _, block := range blocks {
        blockMap[block.ID] = block
    }

    // Создаем новые блоки в указанном порядке
    var newBlocks []models.Block
    for position, blockID := range blockIDs {
        block, exists := blockMap[blockID]
        if !exists {
            return fmt.Errorf("block not found: %s", blockID)
        }
        block.Position = position
        newBlocks = append(newBlocks, block)
    }

    // Сохраняем новый порядок
    return s.blockRepo.SaveBlocks(ctx, materialID, newBlocks)
}