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
