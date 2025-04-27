.PHONY: build run test clean

# 設定変数
APP_NAME=mcp-server-make
BUILD_DIR=./build
CMD_DIR=./cmd/mcp-server-make

# デフォルトターゲット
all: build

# ビルド
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)

# 実行
run:
	@go run $(CMD_DIR)

# テスト実行
test:
	@echo "Running tests..."
	@go test -v ./...

# カバレッジ計測
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# リンターの実行
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

# 依存関係の更新
deps:
	@echo "Updating dependencies..."
	@go mod tidy

# クリーンアップ
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# ヘルプターゲット
help:
	@echo "Available targets:"
	@echo "  build          - アプリケーションをビルド"
	@echo "  run            - アプリケーションを実行"
	@echo "  test           - テストを実行"
	@echo "  test-coverage  - テスト実行とカバレッジレポート作成"
	@echo "  lint           - 静的解析ツールの実行"
	@echo "  deps           - 依存関係の更新"
	@echo "  clean          - ビルドファイルを削除"