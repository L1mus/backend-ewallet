CREATE TABLE newsletter (
    id SERIAL PRIMARY KEY,
    email VARCHAR(254) NOT NULL UNIQUE,
    user_id INT,
    status subscribe_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_newsletter_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE SET NULL
);