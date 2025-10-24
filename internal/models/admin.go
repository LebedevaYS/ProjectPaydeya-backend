package models

import "time"

type AdminStats struct {
    TotalUsers      int `json:"totalUsers"`
    TotalMaterials  int `json:"totalMaterials"`
    ActiveTeachers  int `json:"activeTeachers"`
    PublishedMaterials int `json:"publishedMaterials"`
}

type UserManagement struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    FullName  string    `json:"fullName"`
    Role      string    `json:"role"`
    IsVerified bool     `json:"isVerified"`
    CreatedAt time.Time `json:"createdAt"` // ← Измените на time.Time
    MaterialsCount int  `json:"materialsCount"`
}

type BlockUserRequest struct {
    Reason string `json:"reason" binding:"required"`
}

type CreateSubjectRequest struct {
    ID   string `json:"id" binding:"required"`
    Name string `json:"name" binding:"required"`
    Icon string `json:"icon"`
}