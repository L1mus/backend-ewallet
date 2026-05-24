CREATE TABLE IF NOT EXISTS favorite_contacts (
    id  SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    favorite_user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_transfer_contacts UNIQUE (user_id, favorite_user_id),
    CONSTRAINT fk_tc_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_tc_favorite_user FOREIGN KEY (favorite_user_id)
        REFERENCES users (id) ON DELETE CASCADE
);