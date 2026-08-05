CREATE TABLE IF NOT EXISTS wallet(
    user_id UUID PRIMARY KEY NOT NULL UNIQUE,
    balance_in_cent NUMERIC(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT fk_wallet_user FOREIGN KEY (user_id)
     REFERENCES users (id) ON DELETE CASCADE
);