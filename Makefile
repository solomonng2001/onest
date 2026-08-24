.DEFAULT_GOAL := help

.PHONY: \
	help \
	install \
	dev \
	test \
	backend-run \
	backend-build \
	backend-test \
	backend-check \
	backend-tidy \
	frontend-run \
	frontend-check

help:
	@echo "Available commands:"
	@echo "  make install        Install backend and frontend dependencies"
	@echo "  make dev            Start the backend and frontend"
	@echo "  make test           Check and test the entire project"
	@echo "  make backend-run    Start only the backend"
	@echo "  make backend-build  Build backend/bin/server"
	@echo "  make backend-test   Run backend tests"
	@echo "  make backend-check  Format, vet and test the backend"
	@echo "  make backend-tidy   Clean up Go module dependencies"
	@echo "  make frontend-run   Start only the frontend"
	@echo "  make frontend-check Lint and build the frontend"

install:
	cd backend && go mod download
	cd frontend && npm ci

dev:
	@set -e; \
	(cd backend && go run ./cmd/server) & backend_pid=$$!; \
	(cd frontend && npm run dev) & frontend_pid=$$!; \
	trap 'kill $$backend_pid $$frontend_pid 2>/dev/null || true' INT TERM EXIT; \
	wait

test: backend-check frontend-check

backend-run:
	cd backend && go run ./cmd/server

backend-build:
	mkdir -p backend/bin
	cd backend && go build -o bin/server ./cmd/server

backend-test:
	cd backend && go test ./...

backend-check:
	cd backend && go fmt ./...
	cd backend && go vet ./...
	cd backend && go test ./...

backend-tidy:
	cd backend && go mod tidy

frontend-run:
	cd frontend && npm run dev

frontend-check:
	cd frontend && npm run lint
	cd frontend && npm run build
