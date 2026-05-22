CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    type type_transaction NOT NULL,
    activity_type type_activity_transaction NOT NULL,
    status status_transaction NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transactions_sender FOREIGN KEY (user_id)
        REFERENCES users (id)
);