CREATE TABLE IF NOT EXISTS transfer_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL UNIQUE,
    receiver_id INT NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transfer_details_trx FOREIGN KEY (transaction_id)
        REFERENCES transactions (id) ON DELETE CASCADE,
    CONSTRAINT fk_transfer_details_receiver FOREIGN KEY (receiver_id)
        REFERENCES users (id)
);