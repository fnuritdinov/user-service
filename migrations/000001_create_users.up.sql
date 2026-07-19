CREATE TABLE IF NOT EXISTS users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL,
                       email VARCHAR(100) NOT NULL UNIQUE,
                       phone VARCHAR(20),
                       age INTEGER,
                       password VARCHAR(255) NOT NULL,
                       role TEXT,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);