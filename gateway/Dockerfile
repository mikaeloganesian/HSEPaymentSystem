# syntax=docker/dockerfile:1

FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

ENTRYPOINT ["./main"]
