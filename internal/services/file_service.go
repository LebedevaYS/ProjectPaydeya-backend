package services

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

type FileService struct {
    uploadPath      string
    storageService  *StorageService  // ← ДОБАВИТЬ
}

// Обновленный конструктор
func NewFileService(uploadPath string, storageService *StorageService) *FileService {
    os.MkdirAll(uploadPath, 0755)
    return &FileService{
        uploadPath:     uploadPath,
        storageService: storageService,  // ← ДОБАВИТЬ
    }
}

// UploadImage загружает изображение в облачное хранилище
func (s *FileService) UploadImage(ctx context.Context, file io.Reader, fileName string, userID int) (*UploadResult, error) {
    if s.storageService != nil {
        // Используем облачное хранилище
        return s.storageService.UploadImage(ctx, file, fileName, userID)
    }
    // Fallback на локальное хранилище
    return s.uploadImageLocal(ctx, file, fileName, userID)
}

// UploadVideo загружает видео в облачное хранилище
func (s *FileService) UploadVideo(ctx context.Context, file io.Reader, fileName string, userID int, fileSize int64) (*UploadResult, error) {
    if s.storageService != nil {
        // Используем облачное хранилище
        return s.storageService.UploadVideo(ctx, file, fileName, userID, fileSize)
    }
    // Fallback на локальное хранилище
    return s.uploadVideoLocal(ctx, file, fileName, userID, fileSize)
}

// Локальная загрузка изображения (fallback)
func (s *FileService) uploadImageLocal(ctx context.Context, file io.Reader, fileName string, userID int) (*UploadResult, error) {
    userDir := filepath.Join(s.uploadPath, "images", fmt.Sprintf("%d", userID))
    os.MkdirAll(userDir, 0755)

    ext := filepath.Ext(fileName)
    randomBytes := make([]byte, 8)
    rand.Read(randomBytes)
    newFileName := fmt.Sprintf("%s%s", hex.EncodeToString(randomBytes), ext)
    filePath := filepath.Join(userDir, newFileName)

    out, err := os.Create(filePath)
    if err != nil {
        return nil, err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        return nil, err
    }

    publicURL := fmt.Sprintf("/uploads/images/%d/%s", userID, newFileName)

    return &UploadResult{
        URL:      publicURL,
        FileName: newFileName,
    }, nil
}

// Локальная загрузка видео (fallback)
func (s *FileService) uploadVideoLocal(ctx context.Context, file io.Reader, fileName string, userID int, fileSize int64) (*UploadResult, error) {
    userDir := filepath.Join(s.uploadPath, "videos", fmt.Sprintf("%d", userID))
    os.MkdirAll(userDir, 0755)

    ext := filepath.Ext(fileName)
    randomBytes := make([]byte, 8)
    rand.Read(randomBytes)
    newFileName := fmt.Sprintf("%s%s", hex.EncodeToString(randomBytes), ext)
    filePath := filepath.Join(userDir, newFileName)

    out, err := os.Create(filePath)
    if err != nil {
        return nil, err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        return nil, err
    }

    publicURL := fmt.Sprintf("/uploads/videos/%d/%s", userID, newFileName)

    return &UploadResult{
        URL:      publicURL,
        FileName: newFileName,
        Size:     fileSize,
    }, nil
}

// SaveAvatar сохраняет аватар и возвращает URL (существующий метод)
func (s *FileService) SaveAvatar(userID int, fileHeader *multipart.FileHeader) (string, error) {
    // ... ваш существующий код без изменений ...
    file, err := fileHeader.Open()
    if err != nil {
        return "", fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    // Проверяем тип файла
    buff := make([]byte, 512)
    _, err = file.Read(buff)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %w", err)
    }

    fileType := http.DetectContentType(buff)
    if !strings.HasPrefix(fileType, "image/") {
        return "", fmt.Errorf("file is not an image")
    }

    // Сбрасываем позицию чтения
    file.Seek(0, 0)

    // Генерируем уникальное имя файла
    randomBytes := make([]byte, 8)
    rand.Read(randomBytes)
    fileExt := filepath.Ext(fileHeader.Filename)
    fileName := fmt.Sprintf("avatar_%d_%s%s", userID, hex.EncodeToString(randomBytes), fileExt)
    filePath := filepath.Join(s.uploadPath, "avatars", fileName)

    // Создаем папку avatars если не существует
    os.MkdirAll(filepath.Join(s.uploadPath, "avatars"), 0755)

    // Сохраняем файл
    dst, err := os.Create(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to create file: %w", err)
    }
    defer dst.Close()

    // Копируем содержимое
    if _, err = io.Copy(dst, file); err != nil {
        return "", fmt.Errorf("failed to save file: %w", err)
    }

    // Возвращаем относительный путь
    return "/uploads/avatars/" + fileName, nil
}

// DeleteAvatar удаляет старый аватар (существующий метод)
func (s *FileService) DeleteAvatar(avatarURL string) error {
    if avatarURL == "" {
        return nil
    }

    // Убираем первый слэш для создания правильного пути
    filePath := strings.TrimPrefix(avatarURL, "/")
    return os.Remove(filePath)
}