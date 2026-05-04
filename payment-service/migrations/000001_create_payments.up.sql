CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS payments (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id    UUID        NOT NULL,
    user_id       UUID        NOT NULL,
    amount        DECIMAL(12,2) NOT NULL,
    currency      VARCHAR(3)  NOT NULL DEFAULT 'KZT',
    method        VARCHAR(20) NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    processed_at  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_user_id    ON payments(user_id);
CREATE INDEX idx_payments_status     ON payments(status);
