package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "paydeya-backend/internal/services"
)

func main() {
    // Инициализируем storage service
    storage, err := services.NewStorageService(
        os.Getenv("S3_BUCKET"),
        os.Getenv("S3_ACCESS_KEY"),
        os.Getenv("S3_SECRET_KEY"),
    )
    if err != nil {
        log.Fatalf("❌ Failed to create storage: %v", err)
    }

    fmt.Println("✅ Storage service initialized successfully!")
    fmt.Println("📦 Bucket:", os.Getenv("S3_BUCKET"))
}