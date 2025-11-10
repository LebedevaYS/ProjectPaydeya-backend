package models

import "time"

// AdminStats represents platform statistics for admin
// @Description Статистика платформы для админ-панели
type AdminStats struct {
    TotalUsers      int `json:"totalUsers"`
    TotalMaterials  int `json:"totalMaterials"`
    ActiveTeachers  int `json:"activeTeachers"`
    PublishedMaterials int `json:"publishedMaterials"`
}

// UserManagement represents user data for admin management
// @Description Данные пользователя для управления в админ-панели
type UserManagement struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    FullName  string    `json:"fullName"`
    Role      string    `json:"role"`
    IsVerified bool     `json:"isVerified"`
    IsBlocked bool      `json:"isBlocked" example:"false"`
    BlockReason *string  `json:"blockReason,omitempty" example:"Нарушение правил"`
    CreatedAt time.Time `json:"createdAt"`
    MaterialsCount int  `json:"materialsCount"`
}

// BlockUserRequest represents request to block a user
// @Description Запрос на блокировку пользователя
type BlockUserRequest struct {
    Reason string `json:"reason" binding:"required"`
}

// CreateSubjectRequest represents request to create a new subject
// @Description Запрос на создание нового предмета
type CreateSubjectRequest struct {
    ID   string `json:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
    Icon string `json:"icon"`
}