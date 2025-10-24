package models

type CatalogMaterial struct {
    ID            int     `json:"id"`
    Title         string  `json:"title"`
    Subject       string  `json:"subject"`
    Author        Author  `json:"author"`
    Rating        float64 `json:"rating"`
    StudentsCount int     `json:"studentsCount"`
    Duration      int     `json:"duration,omitempty"`
    Level         string  `json:"level,omitempty"`
    ThumbnailURL  string  `json:"thumbnailUrl,omitempty"`
}

type Author struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Teacher struct {
    ID              int      `json:"id"`
    Name            string   `json:"name"`
    Specializations []string `json:"specializations"`
    Rating          float64  `json:"rating"`
    MaterialsCount  int      `json:"materialsCount"`
    AvatarURL       string   `json:"avatarUrl,omitempty"`
}

type Subject struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Icon string `json:"icon"`
}

type CatalogFilters struct {
    Search  string `form:"search"`
    Subject string `form:"subject"`
    Level   string `form:"level"`
    Page    int    `form:"page"`
    Limit   int    `form:"limit"`
}

type TeacherFilters struct {
    Search  string `form:"search"`
    Subject string `form:"subject"`
}