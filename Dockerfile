FROM golang:1.23-alpine AS builder

#устанавливает рабочую дирректорию /app в контейнере
WORKDIR /app 

#копирует go.mod и go.sum из локальной папки в папку /app в контейнере
COPY go.mod go.sum ./

#скачивает все зависимости
RUN go mod download

#Копирует все файлы из текущей локальной директории в /app контейнера
COPY . . 

#Компилирует Go-приложение из ./cmd/app/ в бинарник /app/bin 
#CGO_ENABLED=0 — статическая компиляция, -ldflags="-s -w" — уменьшает размер бинарника
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin ./cmd/app/

#Создает новый образ на основе чистого Alpine Linux. Уменьшает размер финального образа (удаляются компилятор Go и ненужные инструменты).
FROM alpine:latest

WORKDIR /app

#Устанавливает пакеты tzdata (данные о времени) и ca-certificates (для HTTP-request)
#Флаг --no-cache: Не сохраняет кэш пакетов (уменьшает размер образа).
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /app/bin /app/bin
COPY --from=builder /app/config /app/config 
ENTRYPOINT ["/app/bin"]