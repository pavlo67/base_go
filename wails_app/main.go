package main

import (
	"context"
	"embed"
	"github.com/pavlo67/base_go/lib/logger"
	"github.com/pavlo67/base_go/lib/logger/logger_zap"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"wails-vue-go/backend/appcore"
	"wails-vue-go/backend/modules/filepanel"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Wails needs embedded frontend assets for production builds.
//
//go:embed all:frontend/dist
var assets embed.FS

type App struct {
	ctx        context.Context
	server     *http.Server
	serverOnce sync.Once
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startHTTPServer()
}

func (a *App) shutdown(ctx context.Context) {
	if a.server == nil {
		return
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = a.server.Shutdown(stopCtx)
}

func (a *App) startHTTPServer() {
	a.serverOnce.Do(func() {
		var loggerConfig logger.Config
		l, err := logger_zap.New(loggerConfig)
		if err != nil {
			log.Fatalf("filepanel init failed: %v", err)
		}

		// l := logadapter.New("wails_app")
		registry := appcore.NewRegistry()

		filePanelModule, err := filepanel.NewFromStartup(l)
		if err != nil {
			l.Fatalf("filepanel init failed: %v", err)
		}
		registry.Register(filePanelModule)

		mux := http.NewServeMux()
		registry.RegisterHTTP(mux)

		listener, err := net.Listen("tcp", "127.0.0.1:34116")
		if err != nil {
			l.Fatalf("HTTP backend listen failed: %v", err)
		}

		a.server = &http.Server{
			Handler: appcore.WithLocalDevCORS(mux),
		}

		go func() {
			l.Info("Go HTTP backend: http://127.0.0.1:34116")
			if err := a.server.Serve(listener); err != nil && err != http.ErrServerClosed {
				l.Info("HTTP backend error: %v", err)
			}
		}()
	})
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Wails File Panel",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
	})
	if err != nil {
		log.Fatal(err)
	}
}
