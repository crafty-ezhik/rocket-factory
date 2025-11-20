-- Создаем таблицу с пользователями
create table users (
    user_uuid uuid primary key default gen_random_uuid(),
    login varchar(100) not null ,
    email varchar(100) not null ,
    password text not null ,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Создаем индекс по user_uuid для быстрого поиска пользователей
CREATE INDEX IF NOT EXISTS idx_iam_user_uuid ON users (user_uuid);

-- Создаем таблицу с методами нотификации
create table notification_methods (
    id serial primary key,
    user_uuid UUID not null references users(user_uuid) on delete cascade ,
    provider_name text not null ,
    target text not null
);




