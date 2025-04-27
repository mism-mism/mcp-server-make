package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mism-mism/mcp-server-make/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// ロガーの初期化
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// 終了シグナルを受け取るためのチャネル
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// サーバーの初期化
	srv := server.NewServer()

	// サーバーの起動（非同期）
	go func() {
		if err := srv.Start(":8080"); err != nil {
			log.Fatal().Err(err).Msg("サーバーの起動に失敗しました")
		}
	}()

	log.Info().Msg("サーバーが起動しました（ポート: 8080）")

	// シグナルを待機
	<-sigCh
	log.Info().Msg("シャットダウンシグナルを受信しました")

	// グレースフルシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("サーバーのシャットダウンに問題が発生しました")
	}

	fmt.Println("サーバーを正常に終了しました")
}