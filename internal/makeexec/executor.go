package makeexec

// Executor は make コマンドの実行を担当する構造体
type Executor struct {
	MakePath string
	WorkDir  string
	Timeout  int
}

// NewExecutor は新しい Executor インスタンスを作成します
func NewExecutor(makePath, workDir string, timeout int) *Executor {
	return &Executor{
		MakePath: makePath,
		WorkDir:  workDir,
		Timeout:  timeout,
	}
}
