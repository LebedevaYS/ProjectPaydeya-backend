package main

import (

    "log"
    "os"
    "strconv"

    "paydeya-backend/internal/database"
    "paydeya-backend/internal/handlers"
    "paydeya-backend/internal/repositories"
    "paydeya-backend/internal/services"


    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
 // Загружаем .env файл локально
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  No .env file found, using environment variables")
    }

    // Создаем конфиг для БД
    dbConfig := &database.Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "password"),
        DBName:     getEnv("DB_NAME", "paydeya"),
    }

    // Инициализируем базу данных
    if err := database.Init(dbConfig); err != nil {
        log.Printf("❌ Failed to initialize database: %v", err)
        // НЕ завершаем приложение - возможно мы на Render и БД еще не готова
    } else {
        log.Println("✅ Database connected successfully")
    }
    defer database.Close()

    // Создаем репозитории и сервисы
    userRepo := repositories.NewUserRepository(database.DB)
    authService := services.NewAuthService(userRepo, "your-jwt-secret") // пока используем заглушку
    authHandler := handlers.NewAuthHandler(authService)

    // Настраиваем Gin
    if os.Getenv("GIN_MODE") != "debug" {
        gin.SetMode(gin.ReleaseMode)
    }

    router := gin.Default()

    // CORS middleware
    router.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })

    // Routes
    router.GET("/health", handlers.HealthCheck)
    //router.GET("/api/v1/users", handlers.GetUsersTest(database.DB)) // временный endpoint с подключением к БД
    // Временный эндпоинт для теста
    router.GET("/api/v1/users", func(c *gin.Context) {
        if database.DB == nil {
            c.JSON(500, gin.H{"error": "Database not connected"})
            return
        }

        c.JSON(200, gin.H{
            "message": "Database is connected!",
            "status": "success",
        })
    })

    auth := router.Group("/api/v1/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
        auth.POST("/refresh", authHandler.Refresh)
        auth.POST("/logout", authHandler.Logout)
        auth.POST("/forgot-password", authHandler.ForgotPassword)
        auth.POST("/reset-password", authHandler.ResetPassword)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 Server starting on port %s", port)
    log.Printf("📊 Database connected successfully")
    log.Printf("🌐 Endpoints:")
    log.Printf("   GET /health")
    log.Printf("   GET /api/v1/users")
    log.Printf("   POST /api/v1/auth/register")
    log.Printf("   POST /api/v1/auth/login")
    log.Printf("   POST /api/v1/auth/refresh")
    log.Printf("   POST /api/v1/auth/logout")
    log.Printf("   POST /api/v1/auth/forgot-password")
    log.Printf("   POST /api/v1/auth/reset-password")

    if err := router.Run(":" + port); err != nil {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}

// Вспомогательные функции для env переменных
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}