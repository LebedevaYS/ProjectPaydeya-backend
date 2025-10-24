package models

import "time"

type StudentProgress struct {
    CompletedTopics int     `json:"completedTopics"`
    SuccessRate     float64 `json:"successRate"`
    LearningHours   int     `json:"learningHours"`
    AverageGrade    float64 `json:"averageGrade"`
    CurrentMaterials []ProgressMaterial `json:"currentMaterials"`
}

type ProgressMaterial struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Subject   string    `json:"subject"`
    Progress  float64   `json:"progress"` // 0-100%
    LastActivity time.Time `json:"lastActivity"`
}

type MaterialCompletion struct {
    MaterialID int     `json:"materialId"`
    UserID     int     `json:"userId"`
    TimeSpent  int     `json:"timeSpent"` // в секундах
    Grade      float64 `json:"grade"`     // 1-5
    CompletedAt time.Time `json:"completedAt"`
}

type FavoriteMaterial struct {
    MaterialID int       `json:"materialId"`
    UserID     int       `json:"userId"`
    AddedAt    time.Time `json:"addedAt"`
}