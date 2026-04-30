.PHONY: up down logs build restart clean test swag help

# Default target
help:
	@echo "PCDoma — PC Rental Service"
	@echo ""
	@echo "Usage:"
	@echo "  make up       - Start all services"
	@echo "  make down     - Stop all services"
	@echo "  make build    - Rebuild all images"
	@echo "  make logs     - Follow logs"
	@echo "  make restart  - Restart all services"
	@echo "  make clean    - Remove containers and volumes"
	@echo "  make test     - Run tests in all services"
	@echo "  make swag     - Generate Swagger docs in all services"
	@echo "  make env      - Copy .env.example to .env"

env:
	@if [ ! -f .env ]; then cp .env.example .env && echo ".env created"; else echo ".env already exists"; fi

up: env
	docker-compose up -d

build: env
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

clean:
	docker-compose down -v --remove-orphans

test:
	@for svc in auth-service catalog-service booking-service payment-service notification-service gateway; do \
		echo "=== Testing $$svc ==="; \
		cd $$svc && go test ./... && cd ..; \
	done

swag:
	@for svc in auth-service catalog-service booking-service payment-service notification-service; do \
		echo "=== Generating Swagger for $$svc ==="; \
		cd $$svc && swag init -g cmd/main.go -o docs && cd ..; \
	done

# Individual service logs
logs-auth:
	docker-compose logs -f auth-service

logs-catalog:
	docker-compose logs -f catalog-service

logs-booking:
	docker-compose logs -f booking-service

logs-payment:
	docker-compose logs -f payment-service

logs-notification:
	docker-compose logs -f notification-service

logs-gateway:
	docker-compose logs -f gateway
