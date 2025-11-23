.PHONY: lint, test, sec-scan, build, print-gcl-url, clean-gcl, swagger, githooks

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH_RAW := $(shell uname -m)
ifeq ($(ARCH_RAW),x86_64)
  ARCH := amd64
else ifeq ($(ARCH_RAW),aarch64)
  ARCH := arm64
else
  ARCH := $(ARCH_RAW)
endif

GCL_VERSION ?= v2.4.0
GCL_VER_STR := $(patsubst v%,%,$(GCL_VERSION))
BIN_DIR := tmp
GCL := $(BIN_DIR)/golangci-lint

# URL do tarball
GCL_URL := https://github.com/golangci/golangci-lint/releases/download/$(GCL_VERSION)/golangci-lint-$(GCL_VER_STR)-$(OS)-$(ARCH).tar.gz
GCL_TAR_PATH := golangci-lint-$(GCL_VER_STR)-$(OS)-$(ARCH)/golangci-lint

$(GCL):
	@echo "→ Installing golangci-lint $(GCL_VERSION) for $(OS)/$(ARCH)"
	@mkdir -p $(BIN_DIR)
	@curl -fsSL --retry 3 "$(GCL_URL)" \
	 | tar xz --strip-components=1 -C $(BIN_DIR) "$(GCL_TAR_PATH)"
	@chmod +x $(GCL)
	@echo "✓ golangci-lint installed at $(GCL)"

lint: $(GCL)
	@$(GCL) run -c .golangci.yml ./...

test:
	@go test -cover -covermode=atomic -coverprofile=coverage.out  ./...

sec-scan:
	@command -v govulncheck >/dev/null 2>&1 || { echo "govulncheck not found, installing..."; go install golang.org/x/vuln/cmd/govulncheck@latest; }
	@govulncheck -scan symbol -show verbose ./...

build:
	@go build -o tmp/main ./cmd/api

clean-test:
	@go clean -cache -modcache -i -r
	@rm -f coverage.out

print-gcl-url:
	@echo "URL   : $(GCL_URL)"
	@echo "PATH  : $(GCL_TAR_PATH)"
	@echo "OS/ARCH: $(OS)/$(ARCH)"

clean-gcl:
	@rm -rf $(BIN_DIR)

sqlc:
	@sqlc generate

.PHONY: swagger
swagger:
	@swag init -g internal/config/swagger.go -o api

.PHONY: githooks
githooks:
	@echo "→ Installing git hooks"
	@git config core.hooksPath .githooks
	@echo "✓ Git hooks installed"

# Docker targets
.PHONY: docker-build docker-push docker-run docker-stop docker-clean

docker-build:
	@echo "→ Building production Docker image..."
	@docker build -f Dockerfile.prod -t riskplaceangola/backend-core:latest .
	@echo "✓ Docker image built"

docker-push:
	@echo "→ Pushing Docker image to Docker Hub..."
	@docker push riskplaceangola/backend-core:latest
	@echo "✓ Docker image pushed"

docker-run:
	@echo "→ Running Docker container locally..."
	@docker run -d \
		--name backend_core_local \
		-p 8090:8090 \
		--env-file .env \
		riskplaceangola/backend-core:latest
	@echo "✓ Container started at http://localhost:8090"

docker-stop:
	@echo "→ Stopping Docker container..."
	@docker stop backend_core_local || true
	@docker rm backend_core_local || true
	@echo "✓ Container stopped"

docker-clean:
	@echo "→ Cleaning Docker artifacts..."
	@docker rmi riskplaceangola/backend-core:latest || true
	@docker system prune -f
	@echo "✓ Docker artifacts cleaned"

docker-logs:
	@docker logs -f backend_core_local

# Development helpers
.PHONY: dev-setup dev-run

dev-setup:
	@echo "→ Setting up development environment..."
	@go mod download
	@make githooks
	@echo "✓ Development environment ready"

dev-run:
	@echo "→ Running application in development mode..."
	@go run ./cmd/api