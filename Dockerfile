FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./cmd/app

FROM scratch AS prod
COPY --from=builder /app/cmd/app /cmd/app
CMD ["./cmd/app"]
