\set ON_ERROR_STOP true

CREATE TABLE IF NOT EXISTS app_user (
    email_address VARCHAR(255) PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR NOT NULL,
    public_key VARCHAR NOT NULL
);