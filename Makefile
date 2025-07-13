# Makefile

ENV_FILE := .env

# Load .env file into Make environment (safe for Linux)
include $(ENV_FILE)
export $(shell sed 's/#.*//g' $(ENV_FILE) | grep -E '^\s*[A-Za-z_][A-Za-z0-9_]*=' | xargs)

# Testing for .env loading
echo-env:
	@echo "SERVER_PORT = $(SERVER_PORT)"

# Build Docker containers
build:
	@echo "🚧 Building containers..."
	docker-compose --env-file $(ENV_FILE) build

# Run containers (skip build unless needed)
up:
	@echo "🚀 Starting containers..."
	docker-compose --env-file $(ENV_FILE) up

# Run containers after rebuilding
up-build:
	@echo "🔄 Rebuilding and starting containers..."
	docker-compose --env-file $(ENV_FILE) up --build

# Run containers in detached mode
up-detached:
	@echo "🧠 Starting containers in background..."
	docker-compose--env-file $(ENV_FILE) up -d

# Stop and remove containers
down:
	@echo "🛑 Stopping containers..."
	docker-compose --env-file $(ENV_FILE) down

# View logs
logs:
	docker-compose logs -f --tail=100
