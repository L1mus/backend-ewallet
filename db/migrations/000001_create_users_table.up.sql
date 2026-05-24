CREATE TYPE oauth_provider AS ENUM (
  'google',
  'facebook'
);

CREATE TYPE subscribe_status AS ENUM (
  'active',
  'unsubscribe'
);

CREATE TYPE type_transaction AS ENUM (
  'income',
  'expense'
);

CREATE TYPE type_activity_transaction AS ENUM (
  'transfer',
  'topup'
);

CREATE TYPE status_transaction AS ENUM (
  'pending',
  'success',
  'failed'
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hash_password VARCHAR(255) NOT NULL,
    hash_pin VARCHAR(255),
    phone VARCHAR(16) UNIQUE,
    profile_picture_url VARCHAR(255),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);