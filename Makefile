build:
	@go build -o ./cmd/chatgo

run: build
	@./cmd/chatgo

dev:
	@CompileDaemon -exclude-dir=".git,migrations" \
		-command="./bin/chatgo" \
		-build="go build -o ./bin/chatgo" \
		-color -graceful-kill -log-prefix=false

test:
	@go test -v ./...
