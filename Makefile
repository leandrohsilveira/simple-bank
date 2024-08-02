
DB_PROTOCOL ?= postgresql
DB_USERNAME ?= app
DB_PASSWORD ?= password
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= app
DB_URL ?= $(DB_PROTOCOL)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-up:
	@bash ./migrate --path /database/migrations --database postgresql://app:password@localhost:5432/app?sslmode=disable --verbose up

migrate-down:
	@bash ./migrate --path /database/migrations --database postgresql://app:password@localhost:5432/app?sslmode=disable --verbose down 1

migrate-reset:
	@bash ./migrate --path /database/migrations --database postgresql://app:password@localhost:5432/app?sslmode=disable --verbose down --all

sqlc:
	@sqlc generate

templ:
	@templ generate --path ./server

server-start: sqlc templ
	@go run ./server

server-dev:
	@wgo -dir ./server -dir ./database -xdir ./server/database -file .go -file .templ -file .mod -xfile _templ.go make server-start

client-build:
	@cd ./client && npm run build

client-dev:
	@cd ./client && npm run dev
