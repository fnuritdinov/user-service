CREATE TABLE IF NOT EXISTS hash_tokens (
                             id BIGSERIAL PRIMARY KEY,
                             user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             hash VARCHAR(255) NOT NULL,
                             expired_at TIMESTAMP NOT NULL
);