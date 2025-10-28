-- Удаляем старую таблицу если существует
DROP TABLE IF EXISTS material_ratings;

-- Создаем таблицу рейтингов с правильными foreign keys
CREATE TABLE material_ratings (
    id SERIAL PRIMARY KEY,
    material_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Внешние ключи
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    -- Уникальность: пользователь может оценить материал только один раз
    UNIQUE(material_id, user_id)
);

-- Создаем индекс для быстрого поиска рейтингов по материалу
CREATE INDEX idx_material_ratings_material_id ON material_ratings(material_id);
CREATE INDEX idx_material_ratings_user_id ON material_ratings(user_id);