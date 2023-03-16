build:
	@go build -o ./cmd/chatgo

run: build
	@./cmd/chatgo

dev:
	@CompileDaemon -exclude-dir=".git,migrations" \
		-command="./cmd/chatgo" \
		-build="go build -o ./cmd/chatgo" \
		-color -log-prefix=false

test:
	@go test -v ./...
