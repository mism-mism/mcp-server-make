package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"main/internal/makeexec"
)

var (
	makePath      string
	workdir       string
	timeout       int
	maxConcurrent int64
	debug         bool
)

func init() {
	flag.StringVar(&makePath, "make-path", "make", "Path to make executable")
	flag.StringVar(&workdir, "workdir", ".", "Working directory for make execution")
	flag.IntVar(&timeout, "timeout", 120, "Timeout for make execution in seconds")
	flag.Int64Var(&maxConcurrent, "max-concurrent", 4, "Maximum number of concurrent make executions")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	flag.Parse()
}

func main() {
	// ロギングの設定
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Makefile実行エンジンの作成
	executor := makeexec.NewExecutor(makePath, workdir, timeout, maxConcurrent)

	// シグナルハンドリングの設定
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Info().Msgf("Received signal %s, shutting down...", sig)
		os.Exit(0)
	}()

	// MCPサーバの作成
	s := server.NewMCPServer(
		"MCP Server Make",
		"1.0.0",
		server.WithLogging(),
		server.WithRecovery(),
	)

	// makeツールの定義

	tool := mcp.NewTool("make",
		mcp.WithDescription("Execute make command on a Makefile"),
		mcp.WithString("target",
			mcp.Required(),
			mcp.Description("Make target to execute"),
		),
		mcp.WithString("file",
			mcp.Description("Path to Makefile (optional)"),
		),
		mcp.WithString("workdir",
			mcp.Description("Working directory for make execution (optional)"),
		),
	)

	// ツールハンドラーの登録
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleMakeTool(ctx, request, executor)
	})

	log.Info().
		Str("make-path", makePath).
		Str("workdir", workdir).
		Int("timeout", timeout).
		Int64("max-concurrent", maxConcurrent).
		Msg("MCP Server Make starting...")

	// サーバの起動
	if err := server.ServeStdio(s); err != nil {
		log.Error().Err(err).Msg("Server error")
		os.Exit(1)
	}
}

// handleMakeTool はmakeツールのリクエストを処理するハンドラ関数
func handleMakeTool(ctx context.Context, request mcp.CallToolRequest, executor *makeexec.Executor) (*mcp.CallToolResult, error) {
	log.Debug().Interface("args", request.Params.Arguments).Msg("Make tool called")

	// パラメータの取得
	target, _ := request.Params.Arguments["target"].(string)
	file, _ := request.Params.Arguments["file"].(string)
	workDir, _ := request.Params.Arguments["workdir"].(string)

	// パラメータの検証
	if target == "" {
		return nil, fmt.Errorf("target parameter is required")
	}

	// パラメータをMakeParamsに変換
	params := makeexec.MakeParams{
		Target:  target,
		File:    file,
		WorkDir: workDir,
	}

	// Makefileの実行
	result, err := executor.Execute(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("Make execution failed")
		// エラーが発生しても結果は返す
	}

	// 結果のJSON文字列化
	jsonResult, err := makeexec.SerializeResult(result)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize result: %w", err)
	}

	return mcp.NewToolResultText(jsonResult), nil
}
