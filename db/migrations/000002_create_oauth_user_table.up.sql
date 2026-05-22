CREATE TABLE oauth_user (
    id SERIAL PRIMARY KEY,
    user_id INT  NOT NULL,
    provider_name oauth_provider NOT NULL,
    provider_user_id VARCHAR(255) UNIQUE NOT NULL,
    access_token VARCHAR(255),
    refresh_token VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_oauth_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);