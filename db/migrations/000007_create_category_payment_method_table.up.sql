CREATE TABLE IF NOT EXISTS category_payment_method (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);