package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// グレースフルシャットダウンのタイムアウト時間
const ShutdownTimeout = 10 * time.Second

// Server はHTTPサーバーとルーティングを管理する構造体
type Server struct {
	router chi.Router
	server *http.Server
}

// NewServer は新しいサーバーインスタンスを作成
func NewServer() *Server {
	r := chi.NewRouter()

	// ミドルウェアの設定
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// 基本ルートの設定
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// ルートグループの設定
	r.Route("/api", func(r chi.Router) {
		// API v1 エンドポイント
		r.Route("/v1", func(r chi.Router) {
			// 今後のエンドポイントはここに追加
		})
	})

	return &Server{
		router: r,
	}
}

// Start はHTTPサーバーを指定されたアドレスで起動
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	log.Info().Str("addr", addr).Msg("サーバーを起動します")
	return s.server.ListenAndServe()
}

// Shutdown はサーバーを安全に停止
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	log.Info().Msg("サーバーをシャットダウンしています...")
	return s.server.Shutdown(ctx)
}