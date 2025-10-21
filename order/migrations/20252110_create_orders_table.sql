-- +goose Up
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    part_uuids UUID[] NOT NULL DEFAULT '{}',
    total_price NUMERIC(10,2) NOT NULL CHECK (total_price >= 0),
    transaction_uuid UUID,
    payment_method VARCHAR(30) NOT NULL
        CHECK (payment_method IN ('unknown', 'card', 'sbp', 'credit_card', 'investor_money')),
    status VARCHAR(30) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending_payment', 'paid', 'cancelled')),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE orders
