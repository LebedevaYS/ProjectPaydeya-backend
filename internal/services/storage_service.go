package services

import (
    "context"
    "fmt"
    "io"
    "path/filepath"
    "strings"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/google/uuid"
)

type StorageService struct {
    client *s3.Client
    bucket string
    cdnURL string
}

type UploadResult struct {
    URL      string `json:"url"`
    FileName string `json:"fileName"`
    Size     int64  `json:"size"`
    Width    int    `json:"width,omitempty"`
    Height   int    `json:"height,omitempty"`
}

func NewStorageService(bucket, accessKey, secretKey string) (*StorageService, error) {
    resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
        if service == s3.ServiceID {
            return aws.Endpoint{
                URL:           "https://storage.yandexcloud.net",
                SigningRegion: "ru-central1",
            }, nil
        }
        return aws.Endpoint{}, fmt.Errorf("unknown service: %s", service)
    })

    cfg, err := config.LoadDefaultConfig(context.Background(),
        config.WithEndpointResolverWithOptions(resolver),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
        config.WithRegion("ru-central1"),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    client := s3.NewFromConfig(cfg)

    return &StorageService{
        client: client,
        bucket: bucket,
        cdnURL: fmt.Sprintf("https://%s.storage.yandexcloud.net", bucket),
    }, nil
}

func generateUUID() string {
    return uuid.New().String()
}

func isValidImageExt(ext string) bool {
    supported := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}
    ext = strings.ToLower(ext)
    for _, supportedExt := range supported {
        if ext == supportedExt {
            return true
        }
    }
    return false
}

func getContentType(ext string) string {
    switch strings.ToLower(ext) {
    case ".jpg", ".jpeg":
        return "image/jpeg"
    case ".png":
        return "image/png"
    case ".webp":
        return "image/webp"
    case ".gif":
        return "image/gif"
    case ".mp4":
        return "video/mp4"
    case ".webm":
        return "video/webm"
    default:
        return "application/octet-stream"
    }
}

func (s *StorageService) UploadImage(ctx context.Context, file io.Reader, fileName string, userID int) (*UploadResult, error) {
    ext := strings.ToLower(filepath.Ext(fileName))
    if !isValidImageExt(ext) {
        return nil, fmt.Errorf("unsupported image format: %s. Supported: jpg, jpeg, png, webp, gif", ext)
    }

    newFileName := fmt.Sprintf("images/%d/%s%s", userID, generateUUID(), ext)

    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(newFileName),
        Body:        file,
        ContentType: aws.String(getContentType(ext)),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to upload image: %w", err)
    }

    return &UploadResult{
        URL:      s.cdnURL + "/" + newFileName,
        FileName: newFileName,
    }, nil
}

func (s *StorageService) UploadVideo(ctx context.Context, file io.Reader, fileName string, userID int, fileSize int64) (*UploadResult, error) {
    ext := strings.ToLower(filepath.Ext(fileName))
    if ext != ".mp4" && ext != ".webm" {
        return nil, fmt.Errorf("unsupported video format: %s. Supported: mp4, webm", ext)
    }

    if fileSize > 100*1024*1024 {
        return nil, fmt.Errorf("file too large: %d bytes, maximum is 100MB", fileSize)
    }

    newFileName := fmt.Sprintf("videos/%d/%s%s", userID, generateUUID(), ext)

    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(newFileName),
        Body:        file,
        ContentType: aws.String(getContentType(ext)),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to upload video: %w", err)
    }

    return &UploadResult{
        URL:      s.cdnURL + "/" + newFileName,
        FileName: newFileName,
        Size:     fileSize,
    }, nil
}

func (s *StorageService) DeleteFile(ctx context.Context, fileName string) error {
    _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(s.bucket),
        Key:    aws.String(fileName),
    })
    return err
}