# ============================================
# NATS Framework - Makefile (المصحح)
# ============================================

# ============================================
# المتغيرات
# ============================================
BINARY_NAME := nats
BUILD_DIR := ./dist
CMD_DIR := ./cmd/nats
TEST_DIR := ./tests
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# إصدار التطبيق
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# أعلام البناء
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# الألوان
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
MAGENTA := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

# ============================================
# الأوامر الرئيسية
# ============================================

# الهدف الافتراضي
.DEFAULT_GOAL := help

# ============================================
# أهداف البناء
# ============================================

## build: بناء التطبيق
.PHONY: build
build:
	@echo "$(BLUE)🔨 Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@echo "$(CYAN)📊 Version: $(VERSION)$(NC)"
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "$(CYAN)📊 Size: $$(du -h $(BUILD_DIR)/$(BINARY_NAME) | cut -f1)$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  File not found$(NC)"; \
	fi

## build-linux: بناء التطبيق لنظام Linux
.PHONY: build-linux
build-linux:
	@echo "$(BLUE)🔨 Building for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)/main.go
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64$(NC)"
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ]; then \
		echo "$(CYAN)📊 Size: $$(du -h $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 | cut -f1)$(NC)"; \
	fi

## build-windows: بناء التطبيق لنظام Windows
.PHONY: build-windows
build-windows:
	@echo "$(BLUE)🔨 Building for Windows...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)/main.go
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe$(NC)"
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ]; then \
		echo "$(CYAN)📊 Size: $$(du -h $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe | cut -f1)$(NC)"; \
	fi

## build-macos: بناء التطبيق لنظام macOS
.PHONY: build-macos
build-macos:
	@echo "$(BLUE)🔨 Building for macOS...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)/main.go
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64$(NC)"
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ]; then \
		echo "$(CYAN)📊 Size: $$(du -h $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 | cut -f1)$(NC)"; \
	fi

## build-all: بناء التطبيق لجميع المنصات
.PHONY: build-all
build-all: build-linux build-windows build-macos
	@echo "$(GREEN)✅ All builds complete$(NC)"

## install: تثبيت CLI عالمياً
.PHONY: install
install:
	@echo "$(BLUE)📦 Installing $(BINARY_NAME) CLI...$(NC)"
	@go install $(LDFLAGS) $(CMD_DIR)/main.go
	@echo "$(GREEN)✅ Installed: $(BINARY_NAME)$(NC)"

# ============================================
# أهداف الاختبار
# ============================================

## test: تشغيل الاختبارات
.PHONY: test
test:
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	@go test -v ./...

## test-coverage: تشغيل الاختبارات مع تغطية
.PHONY: test-coverage
test-coverage:
	@echo "$(BLUE)🧪 Running tests with coverage...$(NC)"
	@go test -v -cover -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)✅ Coverage report: $(COVERAGE_HTML)$(NC)"
	@go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print "📊 Total coverage: " $$3}'

## test-unit: تشغيل اختبارات الوحدة
.PHONY: test-unit
test-unit:
	@echo "$(BLUE)🧪 Running unit tests...$(NC)"
	@go test -v ./tests/unit/...

## test-integration: تشغيل اختبارات التكامل
.PHONY: test-integration
test-integration:
	@echo "$(BLUE)🧪 Running integration tests...$(NC)"
	@go test -v ./tests/integration/...

## test-e2e: تشغيل اختبارات شاملة
.PHONY: test-e2e
test-e2e:
	@echo "$(BLUE)🧪 Running E2E tests...$(NC)"
	@go test -v ./tests/e2e/...

## test-all: تشغيل جميع الاختبارات
.PHONY: test-all
test-all: test-unit test-integration test-e2e
	@echo "$(GREEN)✅ All tests passed!$(NC)"

# ============================================
# أهداف التطوير
# ============================================

## run: تشغيل التطبيق في وضع التطوير
.PHONY: run
run:
	@echo "$(BLUE)🚀 Running $(BINARY_NAME)...$(NC)"
	@go run $(CMD_DIR)/main.go serve

## run-debug: تشغيل التطبيق في وضع التصحيح
.PHONY: run-debug
run-debug:
	@echo "$(BLUE)🐛 Running $(BINARY_NAME) in debug mode...$(NC)"
	@go run $(CMD_DIR)/main.go serve --debug

## dev: تشغيل التطبيق مع إعادة تحميل تلقائي (يتطلب air)
.PHONY: dev
dev:
	@echo "$(BLUE)🔄 Running in development mode with hot reload...$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)⚠️  air not installed. Installing...$(NC)"; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

## fmt: تنسيق الكود
.PHONY: fmt
fmt:
	@echo "$(BLUE)📝 Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(NC)"

## lint: فحص الكود (يتطلب golangci-lint)
.PHONY: lint
lint:
	@echo "$(BLUE)🔍 Linting code...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not installed. Skipping...$(NC)"; \
	fi

## vet: فحص الكود
.PHONY: vet
vet:
	@echo "$(BLUE)🔍 Vetting code...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Vetting complete$(NC)"

## tidy: تنظيف الاعتماديات
.PHONY: tidy
tidy:
	@echo "$(BLUE)🧹 Tidying dependencies...$(NC)"
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✅ Dependencies tidied$(NC)"

# ============================================
# أهداف Docker
# ============================================

## docker-build: بناء صورة Docker
.PHONY: docker-build
docker-build:
	@echo "$(BLUE)🐳 Building Docker image...$(NC)"
	@docker build -t nats-framework:$(VERSION) .
	@docker tag nats-framework:$(VERSION) nats-framework:latest
	@echo "$(GREEN)✅ Docker image built: nats-framework:$(VERSION)$(NC)"

## docker-run: تشغيل التطبيق في Docker
.PHONY: docker-run
docker-run:
	@echo "$(BLUE)🐳 Running Docker container...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✅ Container running at http://localhost:8080$(NC)"

## docker-stop: إيقاف تشغيل التطبيق في Docker
.PHONY: docker-stop
docker-stop:
	@echo "$(BLUE)🐳 Stopping Docker container...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✅ Container stopped$(NC)"

## docker-logs: عرض سجلات Docker
.PHONY: docker-logs
docker-logs:
	@docker-compose logs -f

# ============================================
# أهداف التنظيف
# ============================================

## clean: تنظيف الملفات المؤقتة
.PHONY: clean
clean:
	@echo "$(BLUE)🧹 Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE)
	@rm -f $(COVERAGE_HTML)
	@rm -rf storage/logs/*.log
	@rm -rf storage/cache/*
	@rm -rf storage/sessions/*
	@echo "$(GREEN)✅ Clean complete$(NC)"

## clean-all: تنظيف كل شيء
.PHONY: clean-all
clean-all: clean
	@echo "$(BLUE)🧹 Cleaning all...$(NC)"
	@rm -rf storage/*
	@rm -rf vendor/
	@echo "$(GREEN)✅ Clean all complete$(NC)"

# ============================================
# أهداف المساعدة
# ============================================

## help: عرض هذه المساعدة
.PHONY: help
help:
	@echo "$(CYAN)============================================$(NC)"
	@echo "$(CYAN)🚀 NATS Framework - Makefile$(NC)"
	@echo "$(CYAN)============================================$(NC)"
	@echo ""
	@echo "$(YELLOW)📦 Build Commands:$(NC)"
	@echo "  $(GREEN)build$(NC)              - Build the application"
	@echo "  $(GREEN)build-linux$(NC)        - Build for Linux"
	@echo "  $(GREEN)build-windows$(NC)      - Build for Windows"
	@echo "  $(GREEN)build-macos$(NC)        - Build for macOS"
	@echo "  $(GREEN)build-all$(NC)          - Build for all platforms"
	@echo "  $(GREEN)install$(NC)            - Install CLI globally"
	@echo ""
	@echo "$(YELLOW)🧪 Test Commands:$(NC)"
	@echo "  $(GREEN)test$(NC)               - Run all tests"
	@echo "  $(GREEN)test-coverage$(NC)      - Run tests with coverage"
	@echo "  $(GREEN)test-unit$(NC)          - Run unit tests"
	@echo "  $(GREEN)test-integration$(NC)   - Run integration tests"
	@echo "  $(GREEN)test-e2e$(NC)           - Run E2E tests"
	@echo "  $(GREEN)test-all$(NC)           - Run all test suites"
	@echo ""
	@echo "$(YELLOW)🚀 Run Commands:$(NC)"
	@echo "  $(GREEN)run$(NC)                - Run the application"
	@echo "  $(GREEN)run-debug$(NC)          - Run in debug mode"
	@echo "  $(GREEN)dev$(NC)                - Run with hot reload"
	@echo ""
	@echo "$(YELLOW)🔧 Development Commands:$(NC)"
	@echo "  $(GREEN)fmt$(NC)                - Format code"
	@echo "  $(GREEN)lint$(NC)               - Lint code"
	@echo "  $(GREEN)vet$(NC)                - Vet code"
	@echo "  $(GREEN)tidy$(NC)               - Tidy dependencies"
	@echo ""
	@echo "$(YELLOW)🐳 Docker Commands:$(NC)"
	@echo "  $(GREEN)docker-build$(NC)       - Build Docker image"
	@echo "  $(GREEN)docker-run$(NC)         - Run Docker container"
	@echo "  $(GREEN)docker-stop$(NC)        - Stop Docker container"
	@echo "  $(GREEN)docker-logs$(NC)        - View Docker logs"
	@echo ""
	@echo "$(YELLOW)🧹 Clean Commands:$(NC)"
	@echo "  $(GREEN)clean$(NC)              - Clean build artifacts"
	@echo "  $(GREEN)clean-all$(NC)          - Clean everything"
	@echo ""
	@echo "$(CYAN)============================================$(NC)"
	@echo "$(CYAN)📊 Version: $(VERSION)$(NC)"
	@echo "$(CYAN)============================================$(NC)"

# ============================================
# أهداف إضافية
# ============================================

## init: تهيئة مشروع جديد
.PHONY: init
init:
	@echo "$(BLUE)🚀 Initializing new project...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) init

## serve: تشغيل الخادم
.PHONY: serve
serve:
	@echo "$(BLUE)🚀 Starting server...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) serve

## migrate: تشغيل الهجرات
.PHONY: migrate
migrate:
	@echo "$(BLUE)📝 Running migrations...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) migrate

## seed: إضافة بيانات افتراضية
.PHONY: seed
seed:
	@echo "$(BLUE)🌱 Seeding database...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) seed

## deploy: نشر التطبيق
.PHONY: deploy
deploy: build
	@echo "$(BLUE)🚀 Deploying application...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) deploy

# ============================================
# أهداف التثبيت
# ============================================

## setup: إعداد بيئة التطوير
.PHONY: setup
setup: tidy
	@echo "$(BLUE)🔧 Setting up development environment...$(NC)"
	@echo "$(GREEN)✅ Setup complete$(NC)"

## deps: تثبيت الاعتماديات
.PHONY: deps
deps:
	@echo "$(BLUE)📦 Installing dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)✅ Dependencies installed$(NC)"

## tools: تثبيت أدوات التطوير
.PHONY: tools
tools:
	@echo "$(BLUE)🔧 Installing development tools...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/vektra/mockery/v2@latest
	@echo "$(GREEN)✅ Tools installed$(NC)"

# ============================================
# أهداف الإصدار
# ============================================

## release: إنشاء إصدار جديد
.PHONY: release
release: test build-all
	@echo "$(BLUE)📦 Creating release...$(NC)"
	@echo "$(GREEN)✅ Release $(VERSION) created$(NC)"

# ============================================
# مساعدات
# ============================================

## version: عرض الإصدار
.PHONY: version
version:
	@echo "$(CYAN)📊 Version: $(VERSION)$(NC)"
	@echo "$(CYAN)📊 Build Time: $(BUILD_TIME)$(NC)"
	@echo "$(CYAN)📊 Git Commit: $(GIT_COMMIT)$(NC)"

## info: عرض معلومات النظام
.PHONY: info
info:
	@echo "$(CYAN)📊 System Information$(NC)"
	@echo "$(CYAN)============================================$(NC)"
	@echo "  Go Version: $(shell go version)"
	@echo "  OS: $(shell go env GOOS)"
	@echo "  Arch: $(shell go env GOARCH)"
	@echo "  Go Root: $(shell go env GOROOT)"
	@echo "  Go Path: $(shell go env GOPATH)"
	@echo "  Module: $(shell go list -m)"
	@echo "$(CYAN)============================================$(NC)"

# ============================================
# أهداف لمنع الصراعات
# ============================================

# تحويل الأهداف التي تحمل نفس اسم الملفات
.PHONY: FORCE
FORCE:

# تعيين الأهداف الافتراضية
.DEFAULT_GOAL := help