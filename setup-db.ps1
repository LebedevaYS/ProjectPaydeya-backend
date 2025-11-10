
function Write-Info($message) {
    Write-Host "$message" -ForegroundColor Cyan
}

function Write-Success($message) {
    Write-Host "$message" -ForegroundColor Green
}

function Write-Error($message) {
    Write-Host "$message" -ForegroundColor Red
}

Write-Host "Starting database setup with UTF-8..." -ForegroundColor Magenta

# Имя контейнера PostgreSQL
$CONTAINER_NAME = "paydeya-db"

# Проверяем, существует ли база данных через docker exec
Write-Info "Checking database..."
$result = docker exec $CONTAINER_NAME psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='paydeya';"
if ($result.Trim() -eq "1") {
    Write-Success "Database 'paydeya' exists"
} else {
    Write-Error "Database 'paydeya' does not exist in container $CONTAINER_NAME"
    exit 1
}


# Run migrations
Write-Info "Running migrations..."
@(
    "001_create_users_table.sql",
    "002_add_specializations_table.sql", 
    "003_create_materials_tables.sql",
    "004_add_ratings_table.sql",
    "005_create_progress_tables.sql"
) | ForEach-Object {
    Write-Info "Executing $_"
    # Передаем содержимое файла через конвейер в psql
    Get-Content "migrations/$_" | docker exec -i $CONTAINER_NAME psql -U postgres -d paydeya
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to execute $_"
        exit 1
    }
}

# Insert sample data
Write-Info "Inserting sample data..."
Get-Content "migrations/006_sample_data.sql" | docker exec -i $CONTAINER_NAME psql -U postgres -d paydeya
if ($LASTEXITCODE -eq 0) {
    Write-Success "Sample data inserted successfully!"
} else {
    Write-Error "Failed to insert sample data"
    exit 1
}

# Verify data
Write-Info "Verifying data..."
docker exec -i $CONTAINER_NAME psql -U postgres -d paydeya -c "
SELECT 
    'Пользователи: ' || COUNT(*) as users,
    'Материалы: ' || COUNT(*) as materials,
    'Рейтинги: ' || COUNT(*) as ratings
FROM 
    (SELECT 1 FROM users) u,
    (SELECT 1 FROM materials) m,
    (SELECT 1 FROM material_ratings) r;
"

Write-Host "Database setup completed successfully!" -ForegroundColor Green