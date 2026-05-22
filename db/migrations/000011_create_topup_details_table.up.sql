CREATE TABLE topup_details (
    id SERIAL PRIMARY KEY,
    transaction_id    INT NOT NULL UNIQUE,
    payment_method_id INT NOT NULL,
    order_amount NUMERIC(15, 2) NOT NULL,
    delivery_fee NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    tax_amount NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    total_amount NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_topup_trx FOREIGN KEY (transaction_id)
        REFERENCES transactions (id) ON DELETE CASCADE,
    CONSTRAINT fk_topup_payment FOREIGN KEY (payment_method_id)
        REFERENCES payment_method (id)
);