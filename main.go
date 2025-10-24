package main

import (

    "log"
    "os"
    "strconv"

    "paydeya-backend/internal/database"
    "paydeya-backend/internal/handlers"
    "paydeya-backend/internal/repositories"
    "paydeya-backend/internal/services"
    "paydeya-backend/internal/middleware"


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


    // Создаем репозитории и сервисы
    userRepo := repositories.NewUserRepository(database.DB)
    materialRepo := repositories.NewMaterialRepository(database.DB)
    blockRepo := repositories.NewBlockRepository(database.DB)
    authService := services.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
    fileService := services.NewFileService("uploads")
    materialService := services.NewMaterialService(materialRepo, blockRepo)

    // Создаем обработчики
    authHandler := handlers.NewAuthHandler(authService)
    profileHandler := handlers.NewProfileHandler(authService, userRepo, fileService)
    materialHandler := handlers.NewMaterialHandler(materialService)

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

    // Обслуживаем статические файлы (аватары)
    router.Static("/uploads", "./uploads")

    // Routes
    router.GET("/health", handlers.HealthCheck)
    router.GET("/api/v1/users", handlers.GetUsersTest(database.DB)) // временный endpoint с подключением к БД


    auth := router.Group("/api/v1/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
        auth.POST("/refresh", authHandler.Refresh)
        auth.POST("/logout", authHandler.Logout)
        auth.POST("/forgot-password", authHandler.ForgotPassword)
        auth.POST("/reset-password", authHandler.ResetPassword)
    }
    // Защищенные эндпоинты (требуют авторизацию)
    protected := router.Group("/api/v1")
    protected.Use(middleware.AuthMiddleware(authService))
    {
        protected.GET("/profile", profileHandler.GetProfile)
        protected.PATCH("/profile", profileHandler.UpdateProfile)
        protected.POST("/profile/avatar", profileHandler.UploadAvatar)

        protected.POST("/materials", materialHandler.CreateMaterial)
        protected.GET("/materials", materialHandler.GetUserMaterials)
        protected.GET("/materials/:id", materialHandler.GetMaterial)
        protected.PUT("/materials/:id", materialHandler.UpdateMaterial)
        protected.POST("/materials/:id/publish", materialHandler.PublishMaterial)
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
    log.Printf("   GET /api/v1/profile")
    log.Printf("   PATCH /api/v1/profile")
    log.Printf("   POST /api/v1/profile/avatar")
    log.Printf("   POST /api/v1/materials")
    log.Printf("   GET /api/v1/materials")
    log.Printf("   GET /api/v1/materials/:id")
    log.Printf("   PUT /api/v1/materials/:id")
    log.Printf("   POST /api/v1/materials/:id/publish")


    defer func() {
        if database.DB != nil {
            database.Close()
            log.Println("🔌 Database connection closed")
        }
    }()


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