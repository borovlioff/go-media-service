# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Сборка бинаря
RUN go build -o media-server .

# Stage 2: Run
FROM alpine:3.18

# Создаём директорию для файлов
RUN mkdir -p /app/public

WORKDIR /app

# Копируем бинарь из builder
COPY --from=builder /app/media-server .

# Переменные окружения по умолчанию
ENV UPLOAD_DIR=/app/public
ENV PORT=8080
ENV PAGE_LIMIT=20
ENV FILE_MAX_SIZE=104857600
ENV FILE_FILTER=".*"
ENV ALLOWED_MIMES="image/jpeg,image/png,video/mp4"

EXPOSE ${PORT}

CMD ["./media-server"]
