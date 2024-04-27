-- Проверяем существование базы данных db_users
SELECT 'CREATE DATABASE db_users' WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'db_users');

-- Создаем таблицу пользователей в базе данных db_users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL
);

-- Создаем таблицу заметок в базе данных db_users
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(100) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    author VARCHAR(50) NOT NULL
);