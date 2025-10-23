-- migrations/001_create_users_table.sql

-- Таблица пользователей (основная)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('student', 'teacher', 'admin')),
    avatar_url VARCHAR(500),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Таблица специализаций учителей
CREATE TABLE IF NOT EXISTS specializations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    subject VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Вставляем тестовые данные для проверки
INSERT INTO users (email, password_hash, full_name, role, is_verified)
VALUES
    ('student@example.com', 'hashed_password_1', 'Иван Петров', 'student', true),
    ('teacher@example.com', 'hashed_password_2', 'Мария Сидорова', 'teacher', true),
    ('admin@example.com', 'hashed_password_3', 'Администратор Системы', 'admin', true)
ON CONFLICT (email) DO NOTHING;