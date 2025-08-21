# Product Management API Makefile

# Variables
APP_NAME=product-management-api
DOCKER_IMAGE=product-management
DOCKER_TAG=latest
GO_VERSION=1.21
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=product_management

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "${BLUE}$(APP_NAME) - Available commands:${NC}"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "${GREEN}%-20s${NC} %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Setup development environment
	@echo "${BLUE}Setting up development environment...${NC}"
	go mod tidy
	go mod download
	cp .env.example .env
	@echo "${GREEN}Development environment setup complete!${NC}"

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "${BLUE}Installing development tools...${NC}"
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang/mock/mockgen@latest
	@echo "${GREEN}Development tools installed!${NC}"

.PHONY: run
run: ## Run the application
	@echo "${BLUE}Starting $(APP_NAME)...${NC}"
	go run cmd/main.go

.PHONY: build
build: ## Build the application
	@echo "${BLUE}Building $(APP_NAME)...${NC}"
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(APP_NAME) cmd/main.go
	@echo "${GREEN}Build complete! Binary available at bin/$(APP_NAME)${NC}"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "${BLUE}Cleaning build artifacts...${NC}"
	rm -rf bin/
	rm -rf docs/
	go clean
	@echo "${GREEN}Clean complete!${NC}"

.PHONY: test
test: ## Run all tests
	@echo "${BLUE}Running all tests...${NC}"
	go test -v ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "${BLUE}Running unit tests...${NC}"
	go test -v ./test/unit/...

.PHONY: test-integration
test-integration: ## Run integration tests only
	@echo "${BLUE}Running integration tests...${NC}"
	go test -v ./test/integration_test.go

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "${BLUE}Running tests with coverage...${NC}"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report generated at coverage.html${NC}"

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "${BLUE}Generating Swagger documentation...${NC}"
	swag init -g cmd/main.go -o docs
	@echo "${GREEN}Swagger documentation generated!${NC}"

.PHONY: lint
lint: ## Run linter
	@echo "${BLUE}Running linter...${NC}"
	golangci-lint run

.PHONY: format
format: ## Format code
	@echo "${BLUE}Formatting code...${NC}"
	go fmt ./...
	@echo "${GREEN}Code formatted!${NC}"

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "${BLUE}Tidying go modules...${NC}"
	go mod tidy
	@echo "${GREEN}Go modules tidied!${NC}"

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "${BLUE}Building Docker image...${NC}"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "${GREEN}Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)${NC}"

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "${BLUE}Running Docker container...${NC}"
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-compose-up
docker-compose-up: ## Start all services with docker-compose
	@echo "${BLUE}Starting services with docker-compose...${NC}"
	docker-compose up -d
	@echo "${GREEN}Services started! API available at http://localhost:8080${NC}"
	@echo "${GREEN}Swagger UI available at http://localhost:8080/swagger/index.html${NC}"
	@echo "${GREEN}pgAdmin available at http://localhost:5050${NC}"

.PHONY: docker-compose-down
docker-compose-down: ## Stop all services
	@echo "${BLUE}Stopping services...${NC}"
	docker-compose down
	@echo "${GREEN}Services stopped!${NC}"

.PHONY: docker-compose-logs
docker-compose-logs: ## Show docker-compose logs
	docker-compose logs -f

.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "${BLUE}Running database migrations...${NC}"
	go run cmd/main.go migrate
	@echo "${GREEN}Database migrations completed!${NC}"

.PHONY: db-create
db-create: ## Create database
	@echo "${BLUE}Creating database...${NC}"
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
	@echo "${GREEN}Database created!${NC}"

.PHONY: db-drop
db-drop: ## Drop database
	@echo "${RED}Dropping database...${NC}"
	dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
	@echo "${GREEN}Database dropped!${NC}"

.PHONY: db-reset
db-reset: db-drop db-create db-migrate ## Reset database
	@echo "${GREEN}Database reset complete!${NC}"

.PHONY: api-docs
api-docs: swagger ## Open API documentation
	@echo "${BLUE}Opening API documentation...${NC}"
	open http://localhost:8080/swagger/index.html

.PHONY: dev
dev: swagger run ## Start development server with swagger generation

.PHONY: prod-build
prod-build: clean swagger build ## Production build with all steps

.PHONY: check
check: format lint test ## Run all checks (format, lint, test)

.PHONY: pre-commit
pre-commit: format swagger test ## Run pre-commit checks

.PHONY: init-db
init-db: ## Initialize database with sample data
	@echo "${BLUE}Initializing database with sample data...${NC}"
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f scripts/init.sql
	@echo "${GREEN}Database initialized!${NC}"

# Default target
.DEFAULT_GOAL := help