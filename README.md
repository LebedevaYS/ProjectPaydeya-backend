# Paydeya Backend

Бэкенд сервис для проекта Paydeya, написанный на Go.


## 📋 Что нужно установить перед началом

### 1. Установите Git
- **Windows**: Скачайте с [git-scm.com](https://git-scm.com/) → Установите как обычную программу
- **Mac**: Откройте Terminal → Введите `git` → Следуйте инструкциям
- **Linux**: Откройте Terminal → Введите `sudo apt install git`

### 2. Установите Go
- Скачайте с [golang.org/dl/](https://golang.org/dl/)
- Установите как обычную программу
- **Проверьте установку**: Откройте Terminal → Введите `go version` → Должна появиться версия Go

### 3. Установите редактор кода
- **VS Code**
- **Intelliji IDEA**
- Или любой другой редактор

### 4. Установите Docker Desktop
- Windows/Mac: Скачайте с docker.com  
- Установите как обычную программу 
- Запустите Docker Desktop после установки 
- **Проверьте установку**: Откройте Terminal → Введите docker --version

### 5. Запуск базы данных и Redis
```bash
# Запустите контейнеры с БД и Redis
docker compose -f docker-compose.dev.yml up -d

# Проверьте что контейнеры работают
docker ps
```
## 🔐 Настройка окружения

- Скопируйте шаблон настроек
```bash
cp .env.example .env
```
- Отредактируйте .env файл под ваше окружение: заполните настройки БД, JWT секрет и т.д.

## 🚀 Быстрый старт

### Локальная разработка

```bash
# Клонирование репозитория
git clone https://github.com/LebedevaYS/ProjectPaydeya-backend.git
cd ProjectPaydeya-backend

# Установка зависимостей
go mod tidy

# Запуск сервера
go run main.go
```


## 🔧 Git команды (рабочий процесс)

### 🆕 Начало работы над новой фичей

```bash
# Обновить основную ветку
git checkout main
git pull origin main

# Создать новую ветку для фичи
git checkout -b feature/название-фичи
```

### 💾 Ежедневный workflow

```bash
# Проверить статус изменений
git status

# Добавить изменения
git add .

# Создать коммит с понятным сообщением
git commit -m "добавлен новый компонент"

# Отправить изменения
git push origin feature/название-фичи
```

### 🧹 Поддержание чистоты истории

```bash
# Просмотр изменений перед коммитом
git diff

# Отмена непроиндексированных изменений
git restore .

# Отмена последнего коммита (осторожно!)
git reset --soft HEAD~1
```

## 📚 SWAGGER

- Локальная ссылка на Swagger документацию: http://localhost:8080/docs#/auth/post_auth_login

- Общая ссылка: https://paydeya-backend.onrender.com/docs#/auth/post_auth_login