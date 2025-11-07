.PHONY: help frontend backend application clean stop

help:
	@echo "Available commands:"
	@echo "  make frontend     - Run frontend with mocked backend responses"
	@echo "  make backend      - Run backend + database locally (Docker DB, data persists)"
	@echo "  make application  - Run full app with Docker (FE + BE + DB)"
	@echo "  make stop         - Stop all Docker services"
	@echo "  make clean        - Clean build artifacts and Docker volumes"

frontend:
	@trap 'pkill -f "tsx.*mock-server" || true' EXIT; \
	cd frontend && npx tsx ___mocks___/mock-server.ts & \
	sleep 1 && \
	cd frontend && npm run dev

backend:
	@DOCKER_COMPOSE="docker-compose"; \
	if ! command -v docker-compose > /dev/null 2>&1; then \
		DOCKER_COMPOSE="docker compose"; \
	fi; \
	$$DOCKER_COMPOSE -f docker-compose.local.yaml up -d postgres; \
	sleep 2
	cd backend && DB_HOST=localhost DB_PORT=5432 DB_USER=rater DB_PASSWORD=rater DB_NAME=rater DB_SSL_MODE=disable go run main.go

application:
	@DOCKER_COMPOSE="docker-compose"; \
	if ! command -v docker-compose > /dev/null 2>&1; then \
		DOCKER_COMPOSE="docker compose"; \
	fi; \
	$$DOCKER_COMPOSE -f docker-compose.local.yaml up -d postgres backend; \
	sleep 3; \
	echo "Starting frontend..."; \
	cd frontend && npm run dev

stop:
	@DOCKER_COMPOSE="docker-compose"; \
	if ! command -v docker-compose > /dev/null 2>&1; then \
		DOCKER_COMPOSE="docker compose"; \
	fi; \
	$$DOCKER_COMPOSE -f docker-compose.local.yaml down 2>/dev/null || true
	@pkill -f "tsx.*mock-server" || true

clean:
	rm -rf frontend/dist frontend/node_modules/.vite
	@DOCKER_COMPOSE="docker-compose"; \
	if ! command -v docker-compose > /dev/null 2>&1; then \
		DOCKER_COMPOSE="docker compose"; \
	fi; \
	$$DOCKER_COMPOSE -f docker-compose.yaml down -v 2>/dev/null || true

