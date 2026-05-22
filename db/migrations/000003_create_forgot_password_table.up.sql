
CREATE TABLE forgot_password (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_forgot_password_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);