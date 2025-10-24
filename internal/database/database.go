package database

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

/*func Init(cfg *Config) error {
    connString := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
    )

    conn, err := pgx.Connect(context.Background(), connString)
    if err != nil {
        return fmt.Errorf("unable to connect to database: %w", err)
    }

    // Test connection
    if err := conn.Ping(context.Background()); err != nil {
        return fmt.Errorf("unable to ping database: %w", err)
    }

    DB = conn
    log.Println("✅ Successfully connected to PostgreSQL")
    return nil
}
*/
func Init(cfg *Config) error {
    connString := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
    )

    log.Printf("🔗 Connecting to: %s@%s:%d/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

    conn, err := pgx.Connect(context.Background(), connString)
    if err != nil {
        log.Printf("❌ Connection failed: %v", err)
        return fmt.Errorf("unable to connect to database: %w", err)
    }

    // Test connection
    if err := conn.Ping(context.Background()); err != nil {
        log.Printf("❌ Ping failed: %v", err)
        conn.Close(context.Background())
        return fmt.Errorf("unable to ping database: %w", err)
    }

    DB = conn
    log.Printf("✅ Database connected successfully! DB pointer: %p", DB)
    return nil
}

func Close() {
    if DB != nil {
        DB.Close(context.Background())
    }
}

// Config - временно определим здесь, чтобы избежать циклических импортов
type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
}