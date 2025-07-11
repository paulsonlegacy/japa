# Makefile

ENV_FILE := ./internal/config/.env

# Load env variables into Makefile environment (Linux)
include $(ENV_FILE)
export $(shell sed 's/#.*//g' $(ENV_FILE) | grep -E '^\s*[A-Za-z_][A-Za-z0-9_]*=' | xargs)

# Build Docker containers
build:
	docker-compose --env-file $(ENV_FILE) build

# Run containers (after building)
up: build
	docker-compose --env-file $(ENV_FILE) up

# Run containers without building
up-only:
	@echo "server port - $(SERVER_PORT)"
	docker-compose --env-file $(ENV_FILE) up