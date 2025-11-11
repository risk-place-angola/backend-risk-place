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

GCL_VERSION ?= v2.3.0
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
	@govulncheck -scan symbol -show verbose ./...

build:
	@go build -o tmp/main ./cmd/api


print-gcl-url:
	@echo "URL   : $(GCL_URL)"
	@echo "PATH  : $(GCL_TAR_PATH)"
	@echo "OS/ARCH: $(OS)/$(ARCH)"

clean-gcl:
	@rm -rf $(BIN_DIR)

.PHONY: swagger
swagger:
	@swag init -g internal/config/swagger.go -o api

.PHONY: githooks
githooks:
	@echo "→ Installing git hooks"
	@git config core.hooksPath .githooks
	@echo "✓ Git hooks installed"