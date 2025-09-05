CREATE TABLE IF NOT EXISTS users (
    internal_id UUID NOT NULL PRIMARY KEY,
    public_id VARCHAR(510) UNIQUE NOT NULL,
    status VARCHAR(100) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT,
    deleted_at BIGINT
);

CREATE TABLE IF NOT EXISTS user_credentials (
    user_internal_id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(510) NOT NULL,
    phone_number VARCHAR(30),

    CONSTRAINT fk_user_credentials_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id)
);

CREATE TABLE IF NOT EXISTS user_data (
    user_internal_id UUID NOT NULL PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    about VARCHAR(510),
    profile_picture_internal_id UUID,

    CONSTRAINT fk_user_data_user FOREIGN KEY (user_internal_id) REFERENCES users(internal_id)
);