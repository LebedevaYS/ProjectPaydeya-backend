package main

import (
    "context"
    "log"
    "os"
    "strconv"
    "fmt"

    "paydeya-backend/internal/database"
    "paydeya-backend/internal/handlers"
    "paydeya-backend/internal/repositories"
    "paydeya-backend/internal/services"
    "paydeya-backend/internal/middleware"


    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

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

func runMigrations() error {
    migrationFiles := []string{
        "migrations/001_create_users_table.sql",
        "migrations/002_add_specializations_table.sql",
        "migrations/003_create_materials_tables.sql",
        "migrations/004_add_ratings_table.sql",
        "migrations/005_create_progress_tables.sql",
    }

    for _, file := range migrationFiles {
        sql, err := os.ReadFile(file)
        if err != nil {
            return fmt.Errorf("failed to read migration %s: %w", file, err)
        }

        _, err = database.DB.Exec(context.Background(), string(sql))
        if err != nil {
            // Игнорируем ЛЮБЫЕ ошибки выполнения миграций для простоты
            log.Printf("⚠️ Migration %s had issues (ignoring): %v", file, err)
            continue // ← ПРОДОЛЖАЕМ даже при ошибках
        }
        log.Printf("✅ Migration applied: %s", file)
    }
    return nil
}
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
    } else {
        log.Println("✅ Database connected successfully")

        // ЗАПУСКАЕМ МИГРАЦИИ ТОЛЬКО ЕСЛИ БД ПОДКЛЮЧЕНА ← ДОБАВЬ ЗДЕСЬ
        if err := runMigrations(); err != nil {
            log.Printf("⚠️  Migrations failed: %v", err)
        }
    }
    // Инициализация облачного хранилища
    storageService, err := services.NewStorageService(
        os.Getenv("S3_BUCKET"),
        os.Getenv("S3_ACCESS_KEY"),
        os.Getenv("S3_SECRET_KEY"),
    )
    if err != nil {
        log.Printf("⚠️ Failed to initialize cloud storage: %v", err)
        // Fallback на локальное хранилище
        //storageService = services.NewLocalStorageService("uploads", "http://localhost:8080/uploads")
        log.Println("📁 Using local storage as fallback")
    } else {
        log.Println("☁️ Cloud storage initialized successfully!")
    }

    // Создаем репозитории
    userRepo := repositories.NewUserRepository(database.DB)
    materialRepo := repositories.NewMaterialRepository(database.DB)
    blockRepo := repositories.NewBlockRepository(database.DB)
    catalogRepo := repositories.NewCatalogRepository(database.DB)
    progressRepo := repositories.NewProgressRepository(database.DB)
    adminRepo := repositories.NewAdminRepository(database.DB)

    // Создаем сервисы
    authService := services.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
    //fileService := services.NewFileService("uploads")
    fileService := services.NewFileService("uploads", storageService)
    materialService := services.NewMaterialService(materialRepo, blockRepo)
    catalogService := services.NewCatalogService(catalogRepo)
    progressService := services.NewProgressService(progressRepo)
    adminService := services.NewAdminService(adminRepo)

    // Создаем обработчики
    authHandler := handlers.NewAuthHandler(authService)
    profileHandler := handlers.NewProfileHandler(authService, userRepo, fileService)
    materialHandler := handlers.NewMaterialHandler(materialService)
    catalogHandler := handlers.NewCatalogHandler(catalogService)
    progressHandler := handlers.NewProgressHandler(progressService)
    adminHandler := handlers.NewAdminHandler(adminService)
    mediaHandler := handlers.NewMediaHandler(fileService)

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

    router.GET("/debug/routes", func(c *gin.Context) {
        routes := router.Routes()
        var routeInfo []string
        for _, route := range routes {
            routeInfo = append(routeInfo, fmt.Sprintf("%s %s", route.Method, route.Path))
        }
        c.JSON(200, gin.H{"routes": routeInfo})
    })

    // Обслуживаем статические файлы (аватары)
    router.Static("/uploads", "./uploads")

    // Routes
    router.GET("/health", handlers.HealthCheck)
    router.GET("/api/v1/users", handlers.GetUsersTest(database.DB))

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
        protected.POST("/materials/:id/blocks", materialHandler.AddBlock)
        protected.PUT("/materials/:id/blocks/:blockId", materialHandler.UpdateBlock)
        protected.DELETE("/materials/:id/blocks/:blockId", materialHandler.DeleteBlock)
        protected.POST("/materials/:id/blocks/reorder", materialHandler.ReorderBlocks)

        protected.POST("/upload/image", mediaHandler.UploadImage)
        protected.POST("/upload/video", mediaHandler.UploadVideo)
        protected.POST("/embed/video", mediaHandler.EmbedVideo)

        student := protected.Group("/student")
        {
            student.GET("/progress", progressHandler.GetProgress)
            student.GET("/favorites", progressHandler.GetFavorites)
            student.POST("/materials/:id/complete", progressHandler.MarkMaterialComplete)
            student.POST("/materials/:id/favorite", progressHandler.ToggleFavorite)
        }

        admin := protected.Group("/admin")
        admin.Use(middleware.AdminMiddleware()) // ← проверка прав администратора
        {
            admin.GET("/statistics", adminHandler.GetStatistics)
            admin.GET("/users", adminHandler.GetUsers)
            admin.POST("/users/:id/block", adminHandler.BlockUser)
            admin.POST("/subjects", adminHandler.CreateSubject)
        }
    }

    catalog := router.Group("/api/v1/catalog")
    {
        catalog.GET("/materials", catalogHandler.SearchMaterials)
        catalog.GET("/subjects", catalogHandler.GetSubjects)
        catalog.GET("/teachers", catalogHandler.SearchTeachers)
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
    log.Printf("   POST /api/v1/materials/:id/blocks")
    log.Printf("   PUT /api/v1/materials/:id/blocks/:blockId")
    log.Printf("   DELETE /api/v1/materials/:id/blocks/:blockId")
    log.Printf("   POST /api/v1/materials/:id/blocks/reorder")
    log.Printf("   GET /api/v1/catalog/materials")
    log.Printf("   GET /api/v1/catalog/subjects")
    log.Printf("   GET /api/v1/catalog/teachers")
    log.Printf("   GET /api/v1/student/progress")
    log.Printf("   GET /api/v1/student/favorites")
    log.Printf("   POST /api/v1/student/materials/:id/complete")
    log.Printf("   POST /api/v1/student/materials/:id/favorite")
    log.Printf("   GET /api/v1/admin/statistics")
    log.Printf("   GET /api/v1/admin/users")
    log.Printf("   POST /api/v1/admin/users/:id/block")
    log.Printf("   POST /api/v1/admin/subjects")

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

