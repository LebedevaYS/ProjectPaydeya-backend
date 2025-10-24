package models

import (
    "time"
)

type Material struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Subject     string    `json:"subject"`
    AuthorID    int       `json:"authorId"`
    AuthorName  string    `json:"authorName,omitempty"`
    Status      string    `json:"status"` // draft, published, archived
    Access      string    `json:"access"` // open, link
    ShareURL    string    `json:"shareUrl,omitempty"`
    Blocks      []Block   `json:"blocks,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

type Block struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"` // text, image, video, formula, quiz
    Content   map[string]interface{} `json:"content"`
    Styles    map[string]interface{} `json:"styles,omitempty"`
    Position  int                    `json:"position"`
    Animation *BlockAnimation        `json:"animation,omitempty"`
}

type BlockAnimation struct {
    Steps     []AnimationStep `json:"steps,omitempty"`
    Trigger   string          `json:"trigger"` // click, auto
    Delay     int             `json:"delay,omitempty"`
}

type AnimationStep struct {
    Element   string                 `json:"element"`
    Action    string                 `json:"action"` // show, hide, highlight
    Style     map[string]interface{} `json:"style,omitempty"`
}

type CreateMaterialRequest struct {
    Title   string `json:"title" binding:"required"`
    Subject string `json:"subject" binding:"required"`
}

type UpdateMaterialRequest struct {
    Title  string  `json:"title"`
    Blocks []Block `json:"blocks"`
}

type PublishMaterialRequest struct {
    Visibility string `json:"visibility"` // draft, published, archived
    Access     string `json:"access"`     // open, link
}