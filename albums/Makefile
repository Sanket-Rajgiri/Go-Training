
load_env:
	@eval $(cat .env | sed 's/^/export /')

stop_db: 
	@docker-compose down
start_db: 
	@docker-compose up -d 
atlas_install: 
	@curl -sSf https://atlasgo.sh | sh

atlas_diff: 
	@altas migrate diff --env gorm

atlas_apply:
	@atlas migrate apply --env gorm

atlas_status:
	@atlas migrate status --env gorm

run_app:
	@go run .
