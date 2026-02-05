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
