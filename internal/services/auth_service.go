package services

import (
    "context"
    "errors"
    "fmt"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/repositories"

    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo *repositories.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
    }
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
    // Проверяем, существует ли email
    exists, err := s.userRepo.EmailExists(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("error checking email: %w", err)
    }
    if exists {
        return nil, errors.New("user with this email already exists")
    }

    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("error hashing password: %w", err)
    }

    // Создаем пользователя
    user := &models.User{
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FullName:     req.FullName,
        Role:         req.Role,
        IsVerified:   req.Role == "student", // Учеников автоматически верифицируем
    }

    // Сохраняем в БД
    if err := s.userRepo.CreateUser(ctx, user); err != nil {
        return nil, fmt.Errorf("error creating user: %w", err)
    }

    return user, nil
}

// Login выполняет вход пользователя
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.User, error) {
    // Находим пользователя по email
    user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("error finding user: %w", err)
    }
    if user == nil {
        return nil, errors.New("invalid email or password")
    }

    // Проверяем пароль
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid email or password")
    }

    return user, nil
}

// Временно - заглушка для генерации токенов
func (s *AuthService) GenerateTokens(user *models.User) (string, string, error) {
    // TODO: Реализовать JWT токены
    // Пока возвращаем заглушки
    return "access-token-stub", "refresh-token-stub", nil
}