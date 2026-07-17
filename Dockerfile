FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/main.go

FROM scratch

COPY --from=builder /app/app /app

EXPOSE 50051

ENTRYPOINT ["/app"]