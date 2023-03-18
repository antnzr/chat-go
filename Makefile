# Definitions
ROOT                    := $(PWD)
GOLANG_DOCKER_IMAGE     := golang:1.20-alpine
GOLANG_DOCKER_CONTAINER := chatgo

#   Format according to gofmt: https://github.com/cytopia/docker-gofmt
#   Usage:
#       make fmt
#       make fmt path=src/elastic/index_setup.go
fmt:
ifdef path
	@docker run --rm -v ${ROOT}:/data cytopia/gofmt -s -w ${path}
else
	@docker run --rm -v ${ROOT}:/data cytopia/gofmt -s -w .
endif

#   Usage:
#       make lint
lint:
	@docker run --rm -v ${ROOT}:/data -w /data golangci/golangci-lint golangci-lint run

clean:
	@docker rm -f ${GOLANG_DOCKER_CONTAINER} || true

build:
	@go build -o ./cmd/app

run:
	@./cmd/app

start: build run

dev:
	@CompileDaemon -exclude-dir=".git,migrations" \
		-command="./cmd/app" \
		-build="go build -o ./cmd/app" \
		-color -log-prefix=false

test:
	@export ENV=test && go test -v ./...
	# @export ENV=test && go test -mod=mod -count=1 --race ./...
