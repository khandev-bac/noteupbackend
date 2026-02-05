CREATE TABLE IF NOT EXISTS user_coin(
    user_id UUID PRIMARY KEY REFERENCES users(id),
    balance INT NOT NULL DEFAULT 2,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
