CREATE TYPE oauth_provider AS ENUM (
  'google',
  'facebook'
);

CREATE TYPE subscribe_status AS ENUM (
  'active',
  'unsubscribe'
);

CREATE TYPE type_transaction AS ENUM (
  'income',
  'expense'
);

CREATE TYPE type_activity_transaction AS ENUM (
  'transfer',
  'topup'
);

CREATE TYPE status_transaction AS ENUM (
  'pending',
  'success',
  'failed'
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hash_password VARCHAR(255) NOT NULL,
    hash_pin VARCHAR(255),
    phone VARCHAR(16) UNIQUE,
    profile_picture_url VARCHAR(255),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE oauth_user (
    id SERIAL PRIMARY KEY,
    user_id INT  NOT NULL,
    provider_name oauth_provider NOT NULL,
    provider_user_id VARCHAR(255) UNIQUE NOT NULL,
    access_token VARCHAR(255),
    refresh_token VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expired_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_oauth_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

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

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_reviews_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE newsletter (
    id SERIAL PRIMARY KEY,
    email VARCHAR(254) NOT NULL UNIQUE,
    user_id INT,
    status subscribe_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_newsletter_user FOREIGN KEY (user_id)
     REFERENCES users (id) ON DELETE SET NULL
);


CREATE TABLE wallet (
    user_id INT PRIMARY KEY NOT NULL UNIQUE,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_wallet_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE category_payment_method (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

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

CREATE TABLE transfer_details (
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

CREATE TABLE favorite_contacts (
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


CREATE INDEX idx_transactions_user_id ON transactions (user_id);
CREATE INDEX idx_transactions_status ON transactions (status);
CREATE INDEX idx_transactions_created_at ON transactions (created_at DESC);
CREATE INDEX idx_transactions_activity_type ON transactions (activity_type);
CREATE INDEX idx_transfer_details_receiver ON transfer_details (receiver_id);
CREATE INDEX idx_transfer_details_trx_id ON transfer_details (transaction_id);
CREATE INDEX idx_topup_details_trx_id ON topup_details (transaction_id);
CREATE INDEX idx_users_full_name ON users (full_name);
CREATE INDEX idx_users_phone ON users (phone);

SELECT  t.id AS transaction_id, t.amount, t.type, t.activity_type, t.status, t.created_at, u_receiver.full_name AS receiver_name,COUNT(*) OVER() AS total_count
FROM transactions t
LEFT JOIN transfer_details td ON t.id = td.transaction_id
LEFT JOIN users u_receiver ON td.receiver_id = u_receiver.id
WHERE t.user_id = 1
ORDER BY t.created_at DESC
LIMIT 10 OFFSET 0;
