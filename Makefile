# Makefile for MemberLink-Lite

# 变量定义
BINARY_NAME=memberlink-lite
BUILD_DIR=bin
MAIN_PATH=cmd/member-link-lite
CONFIG_DIR=configs

# Go 相关变量
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 构建标签
BUILD_FLAGS=-ldflags="-s -w"

# 默认目标
.PHONY: all
all: clean build

# 构建应用
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

# 构建开发版本（包含调试信息）
.PHONY: build-dev
build-dev:
	@echo "Building development version of $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-dev ./$(MAIN_PATH)

# 构建 Windows 版本
.PHONY: build-windows
build-windows:
	@echo "Building Windows version of $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe ./$(MAIN_PATH)

# 构建 Linux 版本
.PHONY: build-linux
build-linux:
	@echo "Building Linux version of $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./$(MAIN_PATH)

# 构建 macOS 版本
.PHONY: build-macos
build-macos:
	@echo "Building macOS version of $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos ./$(MAIN_PATH)

# 运行应用
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run ./$(MAIN_PATH)

# 运行开发版本
.PHONY: run-dev
run-dev: build-dev
	@echo "Running development version of $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)-dev

# 测试
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 单元测试
.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v ./test/unit/...

# 集成测试
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v ./test/integration/...

# 清理构建文件
.PHONY: clean
clean:
	@echo "Cleaning build files..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# 清理所有生成的文件
.PHONY: clean-all
clean-all: clean
	@echo "Cleaning all generated files..."
	@rm -rf uploads/
	@rm -f *.log
	@rm -f *.sqlite

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v -t -d ./...

# 更新依赖
.PHONY: deps-update
deps-update:
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# 格式化代码
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# 代码检查
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping linting"; \
	fi

# 生成 Swagger 文档
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g $(MAIN_PATH)/main.go -o docs; \
	else \
		echo "swag not found, please install: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# 数据库迁移
.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	@if [ -f "scripts/migrate.sh" ]; then \
		./scripts/migrate.sh; \
	else \
		echo "Migration script not found"; \
	fi

# Docker 构建
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Docker 运行
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 -v $(PWD)/$(CONFIG_DIR):/app/configs $(BINARY_NAME)

# 帮助信息
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-dev      - Build development version"
	@echo "  build-windows  - Build for Windows"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-macos    - Build for macOS"
	@echo "  run            - Run the application"
	@echo "  run-dev        - Run development version"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build files"
	@echo "  clean-all      - Clean all generated files"
	@echo "  deps           - Install dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  swagger        - Generate Swagger docs"
	@echo "  migrate        - Run database migrations"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  help           - Show this help message"