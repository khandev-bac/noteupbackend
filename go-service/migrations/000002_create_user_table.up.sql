CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    password TEXT,
    google_id TEXT,
    user_device TEXT,
    picture TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    plan TEXT NOT NULL DEFAULT 'free' CHECK (plan IN ('free','pro')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
