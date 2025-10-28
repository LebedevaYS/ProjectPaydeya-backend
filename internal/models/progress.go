package models

import "time"

// StudentProgress represents student learning progress
// @Description Прогресс обучения студента
type StudentProgress struct {
    CompletedTopics int                `json:"completedTopics" example:"15"`
    SuccessRate     float64            `json:"successRate" example:"85.5"`
    LearningHours   int                `json:"learningHours" example:"120"`
    AverageGrade    float64            `json:"averageGrade" example:"4.2"`
    CurrentMaterials []ProgressMaterial `json:"currentMaterials"`
}

// ProgressMaterial represents material with progress info
// @Description Материал с информацией о прогрессе
type ProgressMaterial struct {
    ID           int       `json:"id" example:"1"`
    Title        string    `json:"title" example:"Основы алгебры"`
    Subject      string    `json:"subject" example:"math"`
    Progress     float64   `json:"progress" example:"75.5"` // 0-100%
    LastActivity time.Time `json:"lastActivity" example:"2023-01-15T10:30:00Z"`
}

// MaterialCompletion represents material completion record
// @Description Запись о завершении материала
type MaterialCompletion struct {
    MaterialID  int       `json:"materialId" example:"1"`
    UserID      int       `json:"userId" example:"123"`
    TimeSpent   int       `json:"timeSpent" example:"3600"` // в секундах
    Grade       float64   `json:"grade" example:"4.5"`     // 1-5
    CompletedAt time.Time `json:"completedAt" example:"2023-01-15T10:30:00Z"`
}

// FavoriteMaterial represents favorite material record
// @Description Запись об избранном материале
type FavoriteMaterial struct {
    MaterialID int       `json:"materialId" example:"1"`
    UserID     int       `json:"userId" example:"123"`
    AddedAt    time.Time `json:"addedAt" example:"2023-01-15T10:30:00Z"`
}