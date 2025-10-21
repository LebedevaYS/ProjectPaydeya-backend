FROM golang:1.25-alpine

WORKDIR /app

# Сначала копируем только go.mod
COPY go.mod ./

# Скачиваем зависимости (go.sum создастся автоматически если его нет)
RUN go mod download

# Теперь копируем весь остальной код
COPY . .

# Собираем приложение
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]