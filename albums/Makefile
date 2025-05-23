ifneq (,$(wildcard ./.env))
    include .env
    export
endif

stop_db: 
	@docker-compose down
start_db: 
	@docker-compose up -d 
atlas_install: 
	@curl -sSf https://atlasgo.sh | sh

atlas_diff: 
	@atlas migrate diff --env gorm

atlas_apply:
	@atlas migrate apply --env gorm  --url ${DB_URL}

atlas_status:
	@atlas migrate status --env gorm --url ${DB_URL}

run_app:
	@go run .

test_app: 
	@go test ./internal/tests/*

swag_init:
	@swag init -d ./,./internal/handlers
