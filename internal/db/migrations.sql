CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rounds (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    bet_amount NUMERIC(18,2) NOT NULL,
    side VARCHAR(8) NOT NULL,        -- 'heads' | 'tails'
    result VARCHAR(8) NOT NULL,      -- 'heads' | 'tails'
    win_amount NUMERIC(18,2) NOT NULL,
    coefficient NUMERIC(6,3) NOT NULL,
    server_seed VARCHAR(128) NOT NULL,
    client_seed VARCHAR(128) NOT NULL,
    nonce BIGINT NOT NULL,
    server_seed_hash VARCHAR(128) NOT NULL,
    rtp_used NUMERIC(5,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_rounds_user ON rounds(user_id);
CREATE INDEX idx_rounds_created ON rounds(created_at DESC);

CREATE TABLE IF NOT EXISTS rtp_config (
    id INT PRIMARY KEY DEFAULT 1,
    value NUMERIC(5,2) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CHECK (id = 1)
);

INSERT INTO rtp_config (id, value) VALUES (1, 97.0)
ON CONFLICT DO NOTHING;