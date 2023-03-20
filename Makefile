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

lint:
	@docker run --rm -v ${ROOT}:/data -w /data golangci/golangci-lint golangci-lint run

clean:
	@docker rm -f ${GOLANG_DOCKER_CONTAINER} || true

build-docker:
	@docker build -t chatgo .

build:
	@go build -o ./bin/app cmd/chatgo/main.go

run:
	@./bin/app

start: build run

dbmigrate:
	@./scripts/db_migrate.sh migrate

dev:
	@CompileDaemon -exclude-dir=".git,migrations" \
		-command="./bin/app" \
		-build="go build -o ./bin/app cmd/chatgo/main.go" \
		-color -log-prefix=false

test:
	@export ENV=test && go test -v -coverprofile=profile.cov ./...

swagger:
	@swag init -d ./cmd/chatgo,./internal/app/controller --pd -o ./docs --parseInternal --parseDepth 1
