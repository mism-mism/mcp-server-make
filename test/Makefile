
# MCP Server Make テスト用 Makefile
.PHONY: help hello error long-running parallel clean all timeout

# デフォルトターゲット
help:
	@echo "=== MCP Server Make テスト用 Makefile ==="
	@echo "使用可能なターゲット:"
	@echo "  help          このヘルプを表示"
	@echo "  hello         簡単な挨拶を表示"
	@echo "  error         エラーを発生させる"
	@echo "  long-running  長時間実行するターゲット (5秒)"
	@echo "  parallel      複数のタスクを並列実行"
	@echo "  timeout       タイムアウトするターゲット (130秒、デフォルトタイムアウトを超える)"
	@echo "  clean         生成ファイルを削除"
	@echo "  all           hello, parallel, long-runningを順に実行"

# 簡単な挨拶を表示
hello:
	@echo "Hello from MCP Server Make!"
	@echo "Current directory: $(shell pwd)"
	@echo "Current date: $(shell date)"

# エラーを発生させる
error:
	@echo "Generating an error..."
	@exit 1

# 長時間実行するターゲット
long-running:
	@echo "Starting long-running task..."
	@sleep 5
	@echo "Long-running task completed."

# 複数のタスクを並列実行（同時実行数のテスト用）
parallel:
	@echo "Starting parallel tasks..."
	@for i in 1 2 3 4 5; do \
		echo "Task $$i starting..." & \
	done
	@echo "Waiting for all tasks to complete..."
	@wait
	@echo "All parallel tasks completed."

# タイムアウトするターゲット（デフォルトタイムアウト120秒を超える）
timeout:
	@echo "Starting a task that will timeout..."
	@echo "This will run for 130 seconds (default timeout is 120s)"
	@sleep 130
	@echo "This message should not be displayed."

# 生成ファイルを削除
clean:
	@echo "Cleaning up..."
	@rm -f test_output_*.txt
	@echo "Clean completed."

# 複数のターゲットを順番に実行
all: hello parallel long-running
	@echo "All targets completed successfully!"

# 環境変数を表示するターゲット
env:
	@echo "Environment variables:"
	@env | sort

# ファイルを生成するターゲット
generate:
	@echo "Generating test files..."
	@echo "This is test file 1" > test_output_1.txt
	@echo "This is test file 2" > test_output_2.txt
	@echo "Generated files successfully."

# カスタム引数を取得するターゲット（引数: MESSAGE=xxx）
message:
	@if [ -n "$(MESSAGE)" ]; then \
		echo "Custom message: $(MESSAGE)"; \
	else \
		echo "No message provided. Use 'make message MESSAGE=your_message'"; \
	fi