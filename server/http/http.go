package http

import (
	"RestKeyValueStore/config"
	"RestKeyValueStore/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

var srv *http.Server
var shutdownSrv context.CancelFunc

func Startup(mux *http.ServeMux) {
	port := config.Read().Port

	srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux}

	var shutdownCtx context.Context
	shutdownCtx, shutdownSrv = context.WithCancel(context.Background())
	defer shutdownSrv()

	go func() {
		logger.Info(fmt.Sprintf("HTTP Server started up on port: %d", port))

		err := srv.ListenAndServe()

		if err != nil {
			logger.Fatal(err)
		}
	}()

	select {
	case <-shutdownCtx.Done():
		gracefulCtx, cancelGracefulShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelGracefulShutdown()
		if err := srv.Shutdown(gracefulCtx); err != nil {
			logger.Error(fmt.Sprintf("shutdown error: %v", err))
			defer os.Exit(1)
			return
		} else {
			logger.Info("gracefully stopped")
		}
	}
}

func Shutdown() {
	shutdownSrv()
}
