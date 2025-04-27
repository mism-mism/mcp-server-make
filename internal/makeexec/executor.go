package makeexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

// Executor は make コマンドの実行を担当する構造体
type Executor struct {
	MakePath       string
	DefaultWorkDir string
	Timeout        time.Duration
	Semaphore      *semaphore.Weighted
	mu             sync.Mutex
}

// Result は make コマンド実行の結果を表現する構造体
type Result struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
}

// MakeParams はmakeコマンド実行に必要なパラメータを定義
type MakeParams struct {
	Target  string `json:"target"`
	File    string `json:"file,omitempty"`
	WorkDir string `json:"workdir,omitempty"`
}

// NewExecutor は新しい Executor インスタンスを作成します
func NewExecutor(makePath, defaultWorkDir string, timeoutSec int, maxConcurrent int64) *Executor {
	return &Executor{
		MakePath:       makePath,
		DefaultWorkDir: defaultWorkDir,
		Timeout:        time.Duration(timeoutSec) * time.Second,
		Semaphore:      semaphore.NewWeighted(maxConcurrent),
	}
}

// Execute はmakeコマンドを実行し、結果を返します
func (e *Executor) Execute(ctx context.Context, params MakeParams) (*Result, error) {
	// コンテキストとタイムアウトの設定
	if ctx == nil {
		ctx = context.Background()
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()

	// 同時実行数の制御
	if err := e.Semaphore.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("failed to acquire semaphore: %w", err)
	}
	defer e.Semaphore.Release(1)

	// 実行開始時間の記録
	startTime := time.Now()

	// 作業ディレクトリの設定
	workDir := e.DefaultWorkDir
	if params.WorkDir != "" {
		workDir = params.WorkDir
	}

	// コマンド引数の準備
	args := []string{}

	// -f オプションが必要な場合
	if params.File != "" {
		// ファイルパスが相対パスの場合、絶対パスに変換
		if !filepath.IsAbs(params.File) {
			params.File = filepath.Join(workDir, params.File)
		}
		args = append(args, "-f", params.File)
	}

	// ターゲットを追加
	args = append(args, params.Target)

	// コマンドの作成
	cmd := exec.CommandContext(timeoutCtx, e.MakePath, args...)
	cmd.Dir = workDir

	// 標準出力と標準エラー出力のキャプチャ
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// コマンドの実行
	log.Debug().
		Str("make_path", e.MakePath).
		Str("workdir", workDir).
		Strs("args", args).
		Msg("Executing make command")

	err := cmd.Run()

	// 実行時間の計算
	duration := time.Since(startTime)
	durationMs := duration.Milliseconds()

	// 結果の作成
	result := &Result{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		DurationMs: durationMs,
	}

	// エラーハンドリング
	if err != nil {
		// タイムアウトエラーの場合
		if timeoutCtx.Err() == context.DeadlineExceeded {
			result.ExitCode = -1
			result.Error = fmt.Sprintf("make execution timed out after %d seconds", int(e.Timeout.Seconds()))
			return result, fmt.Errorf("make execution timed out: %w", timeoutCtx.Err())
		}

		// 終了コードの取得
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			result.Error = fmt.Sprintf("make exited with code %d", result.ExitCode)
		} else {
			result.ExitCode = -1
			result.Error = fmt.Sprintf("failed to execute make: %v", err)
		}

		return result, fmt.Errorf("make execution failed: %w", err)
	}

	// 成功時は終了コード0
	result.ExitCode = 0
	return result, nil
}

// SerializeResult は Result 構造体を JSON 文字列に変換します
func SerializeResult(result *Result) (string, error) {
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize result: %w", err)
	}
	return string(jsonBytes), nil
}
