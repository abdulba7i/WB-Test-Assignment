# FROM golang:1.24-alpine

# COPY go.mod ./
# COPY go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o api cmd/api/main.go
# RUN go build -o consumer cmd/consumer/main.go
# RUN go build -o publisher cmd/publisher/main.go

FROM golang:1.24-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum, чтобы установить зависимости
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем остальные файлы проекта (включая templates/)
COPY . .

# Собираем бинарники
RUN go build -o api ./cmd/api/main.go
RUN go build -o consumer ./cmd/consumer/main.go
RUN go build -o publisher ./cmd/publisher/main.go
