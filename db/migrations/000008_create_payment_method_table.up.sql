CREATE TABLE payment_method (
    id SERIAL PRIMARY KEY,
    payment_category_id INT NOT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    code VARCHAR(10)  NOT NULL UNIQUE,
    fee  NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    logo_url VARCHAR(255),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payment_method_category FOREIGN KEY (payment_category_id)
        REFERENCES category_payment_method (id)
);