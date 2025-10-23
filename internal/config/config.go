package config

import (
    "os"
    "strconv"
)

type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
    RedisURL   string
    JWTSecret  string
    Port       string
}

func Load() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "password"),
        DBName:     getEnv("DB_NAME", "paydeya"),
        RedisURL:   getEnv("REDIS_URL", "localhost:6379"),
        JWTSecret:  getEnv("JWT_SECRET", "dev-secret-key"),
        Port:       getEnv("PORT", "8080"),
    }
}

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