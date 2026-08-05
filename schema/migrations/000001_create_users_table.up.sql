CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    full_name VARCHAR(156) NOT NULL,
    email VARCHAR(156) UNIQUE NOT NULL,
    hash_password VARCHAR(255) NOT NULL,
    hash_pin VARCHAR(255),
    phone VARCHAR(16) UNIQUE,
    profile_picture_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);