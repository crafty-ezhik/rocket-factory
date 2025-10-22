-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL DEFAULT '{}',
    total_price NUMERIC(10,2) NOT NULL CHECK (total_price >= 0),
    transaction_uuid UUID,
    payment_method VARCHAR(30) NOT NULL
        CHECK (payment_method IN ('UNKNOWN', 'CARD', 'SBP', 'CREDIT_CARD', 'INVESTOR_MONEY')) DEFAULT 'UNKNOWN',
    status VARCHAR(30) NOT NULL
        CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED')) DEFAULT 'PENDING_PAYMENT',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE orders
