#!/bin/bash
# MCP Server Make テスト用スクリプト

# 色の定義
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルパー関数
print_header() {
    echo -e "\n${BLUE}==== $1 ====${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 実行パスの設定
MCP_SERVER="./bin/mcp-server-make"
if [ ! -f "$MCP_SERVER" ]; then
    echo "MCPサーババイナリが見つかりません: $MCP_SERVER"
    echo "先にビルドしてください: go build -o bin/mcp-server-make ./cmd/mcp-server-make"
    exit 1
fi

# MCPリクエストを送信する関数
send_mcp_request() {
    local target=$1
    local file=$2
    local workdir=$3

    # JSONリクエストの作成
    local json_request='{
        "jsonrpc": "2.0",
        "id": 1,
        "method": "callTool",
        "params": {
            "tool": "make",
            "arguments": {
                "target": "'$target'"'

    # ファイルが指定されている場合は追加
    if [ -n "$file" ]; then
        json_request+=',
                "file": "'$file'"'
    fi

    # 作業ディレクトリが指定されている場合は追加
    if [ -n "$workdir" ]; then
        json_request+=',
                "workdir": "'$workdir'"'
    fi

    # JSONリクエストを閉じる
    json_request+='
            }
        }
    }'

    # デバッグ用にリクエストを表示
    echo "送信リクエスト:"
    echo "$json_request" | jq '.'

    # MCPサーバにリクエストを送信
    echo "$json_request" | $MCP_SERVER --debug

    echo ""
}

# テストケース
print_header "基本的なhelloターゲットのテスト"
send_mcp_request "hello"

print_header "エラーを生成するターゲットのテスト"
send_mcp_request "error"

print_header "長時間実行するターゲットのテスト"
send_mcp_request "long-running"

print_header "カスタムMakefileパスを指定してテスト"
send_mcp_request "hello" "Makefile"

print_header "カスタム作業ディレクトリを指定してテスト"
send_mcp_request "hello" "" "./"

print_header "並列実行のテスト"
send_mcp_request "parallel"

print_header "環境変数表示のテスト"
send_mcp_request "env"

print_header "ファイル生成のテスト"
send_mcp_request "generate"
if [ -f "test_output_1.txt" ] && [ -f "test_output_2.txt" ]; then
    print_success "ファイルが正常に生成されました"
else
    print_error "ファイル生成に失敗しました"
fi

print_header "クリーンアップのテスト"
send_mcp_request "clean"
if [ ! -f "test_output_1.txt" ] && [ ! -f "test_output_2.txt" ]; then
    print_success "ファイルが正常に削除されました"
else
    print_error "ファイル削除に失敗しました"
fi

print_header "カスタムメッセージのテスト"
export MESSAGE="Hello from test script"
send_mcp_request "message"

# 注: タイムアウトテストは長時間かかるのでコメントアウト
# print_header "タイムアウトのテスト"
# send_mcp_request "timeout"

print_header "全テスト完了"
echo "テストが完了しました。"