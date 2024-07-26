
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