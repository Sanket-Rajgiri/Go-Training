ifneq (,$(wildcard ./.env.list))
    include ./.env.list
    export
endif


stop_db: 
	@docker-compose down
start_db: 
	@docker-compose up -d mysql-db
atlas_install: 
	@curl -sSf https://atlasgo.sh | sh

atlas_diff: 
	@atlas migrate diff --env gorm

atlas_apply:
	@atlas migrate apply --env gorm  --url ${DB_URL}

atlas_status:
	@atlas migrate status --env gorm --url ${DB_URL}

build_app:
	@go build -o albums-app ./cmd/albums

build_docker_app:
	@docker-compose build gin-app
run_docker_app: 
	@docker-compose up -d gin-app
run_app:
	@go run ./cmd/albums

test_app: 
	@go test ./internal/tests/*

swag_init:
	@swag init -d ./cmd/albums,./internal/handlers
