# Этап 1: Сборка бинарного файла приложения
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache git

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код проекта
COPY . .

# Компилируем приложение в один бинарный файл с отключением CGO (для работы в чистом alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o /sub-aggregator ./cmd/app/main.go

# Этап 2: Финальный минималистичный образ для выполнения
FROM alpine:3.21

WORKDIR /root/

# Копируем скомпилированный бинарник из первого этапа
COPY --from=builder /sub-aggregator .
# Копируем папки миграций для утилиты migrate
COPY --from=builder /app/migrations ./migrations

# Открываем порт наружу
EXPOSE 8080

# Команда запуска по умолчанию
CMD ["./sub-aggregator"]
