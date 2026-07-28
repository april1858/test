-- Включаем расширение для генерации UUID, если оно понадобится в будущем
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS subscriptions (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
service_name VARCHAR(255) NOT NULL,
price INT NOT NULL, -- По ТЗ стоимость — целое число рублей
user_id UUID NOT NULL,
start_date DATE NOT NULL, -- Храним как DATE (всегда 1-е число месяца) для удобства индексации
end_date DATE, -- Опциональное поле окончания подписки
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индексы для ускорения работы агрегационной ручки (Поиск по пользователю и сервису за период)
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_service ON subscriptions(user_id, service_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_dates ON subscriptions(start_date, end_date);