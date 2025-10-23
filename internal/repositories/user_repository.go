package repositories

import (
    "context"

    "paydeya-backend/internal/models"

    "github.com/jackc/pgx/v5"
)

type UserRepository struct {
    db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
    return &UserRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (email, password_hash, full_name, role, avatar_url, is_verified)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

    err := r.db.QueryRow(ctx, query,
        user.Email, user.PasswordHash, user.FullName, user.Role, user.AvatarURL, user.IsVerified,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

    return err
}

// GetUserByEmail возвращает пользователя по email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User

    query := `
        SELECT id, email, password_hash, full_name, role, avatar_url, is_verified, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    err := r.db.QueryRow(ctx, query, email).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role,
        &user.AvatarURL, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
    )

    if err == pgx.ErrNoRows {
        return nil, nil
    }

    return &user, err
}

// GetUserByID возвращает пользователя по ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    var user models.User

    query := `
        SELECT id, email, password_hash, full_name, role, avatar_url, is_verified, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    err := r.db.QueryRow(ctx, query, id).Scan(
        &user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role,
        &user.AvatarURL, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt,
    )

    if err == pgx.ErrNoRows {
        return nil, nil
    }

    return &user, err
}

// EmailExists проверяет, существует ли email
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    err := r.db.QueryRow(ctx, query, email).Scan(&exists)
    return exists, err
}