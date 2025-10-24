#!/bin/bash
echo "🚀 Starting database migration on Render..."

# Используем DATABASE_URL которая автоматически есть на Render
psql $DATABASE_URL -f migrations/001_create_users_table.sql

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
else
    echo "❌ Migration failed!"
    exit 1
fi