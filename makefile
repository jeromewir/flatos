.PHONY: all build clean dev install docker-build docker-run

all: build

build:
	go build -o flatos main.go

clean:
	rm -f flatos

dev:
	go tool github.com/air-verse/air

install:
	go mod download

docker-build:
	docker build -t flatos-api:latest .

docker-run:
	docker run -it -p 8080:8080 -v ./database:/app/database flatos-api:latest

install-deps:
# We need this step as go tool does not support -tags flag 
	go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate:
# When updating this command, also make sure to update the Dockerfile
# This command is used in dev
	migrate -path ./sql/migrations -database sqlite://database/flatos.sqlite up

generate-sqlc:
	go tool sqlc generate