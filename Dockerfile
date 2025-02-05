FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/app/main.go


FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
CMD ["./server"]

