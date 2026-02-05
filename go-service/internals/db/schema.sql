CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

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


CREATE TABLE IF NOT EXISTS notes(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,
    audio_url TEXT,
    audio_duration_seconds INT ,
    audio_file_size_mb INT,
    transcript TEXT,
    title TEXT,
    word_count INT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','processing','completed')),
    search_vector TSVECTOR GENERATED ALWAYS AS (
          setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
          setweight(to_tsvector('english', coalesce(transcript, '')), 'B')
    ) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Full-text search (BEST)
CREATE INDEX IF NOT EXISTS notes_search_vector_idx
ON notes USING GIN (search_vector);

-- Title substring search (ILIKE %text%)
CREATE INDEX IF NOT EXISTS notes_title_trgm_idx
ON notes USING GIN (title gin_trgm_ops);

-- User lookup
CREATE INDEX IF NOT EXISTS notes_user_id_idx
ON notes (user_id);

-- Status filtering
CREATE INDEX IF NOT EXISTS notes_status_idx
ON notes (status);



CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    revenucat_subscription_id TEXT,
    revenucat_customer_id TEXT,

    plan TEXT NOT NULL DEFAULT 'free'
        CHECK (plan IN ('free', 'pro')),

    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'cancelled', 'expired', 'paused', 'pending')),

    start_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    renewal_date TIMESTAMPTZ NOT NULL,
    expiration_date TIMESTAMPTZ NOT NULL,
    cancellation_date TIMESTAMPTZ,
    cancellation_reason TEXT,

    billing_cycle_days INT NOT NULL DEFAULT 30,
    amount_paid INT NOT NULL DEFAULT 0,
    currency TEXT,

    original_transaction_id TEXT,
    last_transaction_id TEXT,

    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);



-- Foreign key lookup
CREATE INDEX subscriptions_user_id_idx
ON subscriptions (user_id);

-- Status filtering
CREATE INDEX subscriptions_status_idx
ON subscriptions (status);

-- Plan filtering
CREATE INDEX subscriptions_plan_idx
ON subscriptions (plan);

-- Active subscriptions lookup
CREATE INDEX subscriptions_active_idx
ON subscriptions (user_id)
WHERE status = 'active' AND is_deleted = FALSE;


CREATE TABLE IF NOT EXISTS coin_packs(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    coin_value INT  NOT NULL,
    coin_price INT NOT NULL,
    popular BOOLEAN DEFAULT FALSE,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX coin_packs_idx
ON coin_packs(id);

CREATE TABLE IF NOT EXISTS user_coin(
    user_id UUID PRIMARY KEY REFERENCES users(id),
    balance INT NOT NULL DEFAULT 2,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS coin_transactions(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    amount INT NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX coin_transactions_idx
ON coin_transactions(id);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER notes_updated_at
BEFORE UPDATE ON notes
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER subscriptions_updated_at
BEFORE UPDATE ON subscriptions
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
