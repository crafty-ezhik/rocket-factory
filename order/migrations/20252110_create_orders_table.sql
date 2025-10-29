-- +goose Up

-- создаем таблицу заказов
CREATE TABLE orders (
    order_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    part_uuids UUID[] NOT NULL DEFAULT '{}',
    total_price NUMERIC(10,2) NOT NULL CHECK (total_price >= 0),
    transaction_uuid UUID,
    payment_method VARCHAR(30) NOT NULL DEFAULT 'UNKNOWN',
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING_PAYMENT',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE
);

-- создаем индекс по user_uuid для быстрого поиска заказов пользователя
CREATE INDEX IF NOT EXISTS idx_orders_user_uuid ON orders (user_uuid);


-- +goose Down
DROP TABLE orders

-- удаляем таблицу заказов
DROP TABLE IF EXISTS orders;

-- удаляем индекс по user_uuid
DROP INDEX IF EXISTS idx_orders_user_uuid;