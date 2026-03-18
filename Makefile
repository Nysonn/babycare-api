.PHONY: dev down logs clean migrate-up migrate-down

## dev: build and start all services with hot reload
dev:
	docker-compose up --build

## down: stop all running services
down:
	docker-compose down

## logs: tail live logs from the api service
logs:
	docker-compose logs -f api

## clean: stop services, remove volumes and orphan containers
clean:
	docker-compose down -v --remove-orphans

## migrate-up: run all pending database migrations
migrate-up:
	goose -dir db/migrations postgres $$(grep DATABASE_URL .env | cut -d '=' -f2) up

## migrate-down: roll back the last database migration
migrate-down:
	goose -dir db/migrations postgres $$(grep DATABASE_URL .env | cut -d '=' -f2) down
