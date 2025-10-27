-- Таблица рейтингов материалов
CREATE TABLE IF NOT EXISTS material_ratings (
    id SERIAL PRIMARY KEY,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(material_id, user_id)
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_material_ratings_material_id ON material_ratings(material_id);
CREATE INDEX IF NOT EXISTS idx_material_ratings_user_id ON material_ratings(user_id);
