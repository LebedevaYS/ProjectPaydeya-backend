package models

import (
    "time"
)

// Material represents educational material
// @Description Учебный материал
type Material struct {
    ID          int       `json:"id" example:"1"`
    Title       string    `json:"title" example:"Основы алгебры"`
    Subject     string    `json:"subject" example:"math"`
    AuthorID    int       `json:"authorId" example:"123"`
    AuthorName  string    `json:"authorName,omitempty" example:"Иван Иванов"`
    Status      string    `json:"status" example:"published"` // draft, published, archived
    Access      string    `json:"access" example:"open"` // open, link
    ShareURL    string    `json:"shareUrl,omitempty" example:"https://paydeya.com/share/abc123"`
    Blocks      []Block   `json:"blocks,omitempty"`
    CreatedAt   time.Time `json:"createdAt" example:"2023-01-15T10:30:00Z"`
    UpdatedAt   time.Time `json:"updatedAt" example:"2023-01-15T10:30:00Z"`
}

// Block represents content block in material
// @Description Блок контента в материале
type Block struct {
    ID        string                 `json:"id" example:"block_123"`
    Type      string                 `json:"type" example:"text"` // text, image, video, formula, quiz
    Content   map[string]interface{} `json:"content"`
    Styles    map[string]interface{} `json:"styles,omitempty"`
    Position  int                    `json:"position" example:"1"`
    Animation *BlockAnimation        `json:"animation,omitempty"`
}

// BlockAnimation represents block animation
// @Description Анимация блока
type BlockAnimation struct {
    Steps     []AnimationStep `json:"steps,omitempty"`
    Trigger   string          `json:"trigger" example:"click"` // click, auto
    Delay     int             `json:"delay,omitempty" example:"1000"`
}

// AnimationStep represents animation step
// @Description Шаг анимации
type AnimationStep struct {
    Element   string                 `json:"element" example:"element_1"`
    Action    string                 `json:"action" example:"show"` // show, hide, highlight
    Style     map[string]interface{} `json:"style,omitempty"`
}

// CreateMaterialRequest represents create material request
// @Description Запрос на создание материала
type CreateMaterialRequest struct {
    Title   string `json:"title" binding:"required" example:"Новый материал"`
    Subject string `json:"subject" binding:"required" example:"math"`
}

// UpdateMaterialRequest represents update material request
// @Description Запрос на обновление материала
type UpdateMaterialRequest struct {
    Title  string  `json:"title" example:"Обновленное название"`
    Blocks []Block `json:"blocks"`
}

// PublishMaterialRequest represents publish material request
// @Description Запрос на публикацию материала
type PublishMaterialRequest struct {
    Visibility string `json:"visibility" example:"published"` // draft, published, archived
    Access     string `json:"access" example:"open"`     // open, link
}