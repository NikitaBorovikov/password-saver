CREATE TABLE IF NOT EXISTS passwords(
    password_id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL,
    enc_service TEXT NOT NULL,
    enc_password TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);