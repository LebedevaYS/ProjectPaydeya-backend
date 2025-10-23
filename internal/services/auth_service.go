package services

import (
    "context"
    "errors"
    "fmt"
    "strconv"

    "paydeya-backend/internal/models"
    "paydeya-backend/internal/repositories"
    "paydeya-backend/internal/utils"

    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
    userRepo  *repositories.UserRepository
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

// GenerateTokens создает access и refresh токены
func (s *AuthService) GenerateTokens(user *models.User) (string, string, error) {
    accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, s.jwtSecret)
    if err != nil {
        return "", "", fmt.Errorf("error generating access token: %w", err)
    }

    refreshToken, err := utils.GenerateRefreshToken(user.ID, s.jwtSecret)
    if err != nil {
        return "", "", fmt.Errorf("error generating refresh token: %w", err)
    }

    return accessToken, refreshToken, nil
}

// ValidateToken проверяет токен
func (s *AuthService) ValidateToken(tokenString string) (*utils.Claims, error) {
    return utils.ValidateToken(tokenString, s.jwtSecret)
}

// RefreshTokens обновляет токены
func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
    token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })

    if err != nil || !token.Valid {
        return "", "", errors.New("invalid refresh token")
    }

    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok {
        return "", "", errors.New("invalid token claims")
    }

    // Получаем userID из subject
    userIDStr := claims.Subject
    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        return "", "", errors.New("invalid user ID in token")
    }

    // Находим пользователя
    user, err := s.userRepo.GetUserByID(context.Background(), userID)
    if err != nil || user == nil {
        return "", "", errors.New("user not found")
    }

    // Генерируем новые токены
    return s.GenerateTokens(user)
}

// ForgotPassword - отправка email с токеном сброса
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
    // Пока заглушка - в реальности отправили бы email
    fmt.Printf("📧 Password reset requested for: %s\n", email)
    fmt.Printf("🔗 Reset token: reset-token-%s\n", email) // временный токен
    return nil
}

// ResetPassword - сброс пароля по токену
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
    // Пока заглушка - в реальности проверили бы токен в БД
    fmt.Printf("🔄 Password reset with token: %s\n", token)

    // Хешируем новый пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("error hashing password: %w", err)
    }

    fmt.Printf("✅ New password hash: %s\n", string(hashedPassword))
    return nil
}