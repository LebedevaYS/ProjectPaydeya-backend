-- Таблица завершения материалов
CREATE TABLE IF NOT EXISTS material_completions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    time_spent INTEGER NOT NULL DEFAULT 0, -- в секундах
    grade DECIMAL(3,2) CHECK (grade >= 1 AND grade <= 5),
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, material_id)
);

-- Таблица избранных материалов
CREATE TABLE IF NOT EXISTS favorite_materials (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, material_id)
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_material_completions_user_id ON material_completions(user_id);
CREATE INDEX IF NOT EXISTS idx_material_completions_material_id ON material_completions(material_id);
CREATE INDEX IF NOT EXISTS idx_favorite_materials_user_id ON favorite_materials(user_id);

-- Тестовые данные
INSERT INTO material_completions (user_id, material_id, time_spent, grade) VALUES
(4, 1, 3600, 4.5),
(4, 2, 1800, 5.0)
ON CONFLICT (user_id, material_id) DO NOTHING;

INSERT INTO favorite_materials (user_id, material_id) VALUES
(4, 1),
(4, 3)
ON CONFLICT (user_id, material_id) DO NOTHING;