CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS bookings (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID         NOT NULL,
    booking_type     VARCHAR(10)  NOT NULL,
    pc_id            VARCHAR(24),
    setup_id         VARCHAR(24),
    location_id      VARCHAR(50),
    worker_id        UUID,
    rental_type      VARCHAR(10)  NOT NULL,
    start_time       TIMESTAMPTZ  NOT NULL,
    end_time         TIMESTAMPTZ  NOT NULL,
    base_price       DECIMAL(12,2) NOT NULL,
    extras_price     DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount         DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_price      DECIMAL(12,2) NOT NULL,
    currency         VARCHAR(3)   NOT NULL DEFAULT 'KZT',
    status           VARCHAR(30)  NOT NULL DEFAULT 'pending_payment',
    rejection_reason TEXT,
    payment_id       UUID,
    accepted_at      TIMESTAMPTZ,
    delivered_at     TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS booking_peripherals (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id    UUID         NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    peripheral_id VARCHAR(24)  NOT NULL,
    price_per_unit DECIMAL(12,2) NOT NULL DEFAULT 0
);

CREATE INDEX idx_bookings_user_id     ON bookings(user_id);
CREATE INDEX idx_bookings_status      ON bookings(status);
CREATE INDEX idx_bookings_location_id ON bookings(location_id);
CREATE INDEX idx_bp_booking_id        ON booking_peripherals(booking_id);
