-- Удаляем таблицу пользователей
drop table if exists users;

-- удаляем таблицу методов уведомлений
drop table if exists notification_methods;

-- удаляем индекс
drop index if exists idx_iam_user_uuid;