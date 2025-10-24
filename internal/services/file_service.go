package services

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "strings"
    "net/http"

)

type FileService struct {
    uploadPath string
}

func NewFileService(uploadPath string) *FileService {
    // Создаем папку если не существует
    os.MkdirAll(uploadPath, 0755)
    return &FileService{uploadPath: uploadPath}
}

// SaveAvatar сохраняет аватар и возвращает URL
func (s *FileService) SaveAvatar(userID int, fileHeader *multipart.FileHeader) (string, error) {
    // Открываем файл
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

// DeleteAvatar удаляет старый аватар
func (s *FileService) DeleteAvatar(avatarURL string) error {
    if avatarURL == "" {
        return nil
    }

    // Убираем первый слэш для создания правильного пути
    filePath := strings.TrimPrefix(avatarURL, "/")
    return os.Remove(filePath)
}