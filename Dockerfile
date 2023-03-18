FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./bin/app cmd/chatgo/main.go

FROM scratch AS prod
COPY --from=builder /app/bin/app /bin/app
CMD ["./bin/app"]
