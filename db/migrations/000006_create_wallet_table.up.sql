CREATE TABLE wallet (
    user_id INT PRIMARY KEY NOT NULL UNIQUE,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_wallet_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);