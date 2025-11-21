-- удаляем таблицу заказов
DROP TABLE IF EXISTS orders;

-- удаляем индекс по user_uuid
DROP INDEX IF EXISTS idx_orders_user_uuid;