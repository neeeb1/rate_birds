include .env
export

BINARY_NAME=go-app

build:
	go build -o ./bin/${BINARY_NAME} main.go

run:
	go run

clean:
	go clean
	rm ./bin/${BINARY_NAME}

compose-up:
	docker compose up --build -d

compose-down:
	docker compose down

sqlc:
	sqlc generate

goose-up:
	goose -dir ./sql/schema/ postgres "$(DB_URL)" up

goose-down:
	goose -dir ./sql/schema/ postgres "$(DB_URL") down
