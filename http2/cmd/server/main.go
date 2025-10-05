package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"go-code-patterns/http2/pkg/handlers"
	"go-code-patterns/http2/web"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var port int
	var isHttp1 bool
	var http1FullDuplex bool
	var flushResp bool
	var certPath string
	var keyPath string
	flag.IntVar(&port, "port", 8080, "")
	flag.BoolVar(&isHttp1, "http1", false, "")
	flag.BoolVar(&http1FullDuplex, "http1-full-duplex", false, "")
	flag.BoolVar(&flushResp, "flush", true, "")
	flag.StringVar(&certPath, "cert", "configs/cert.pem", "")
	flag.StringVar(&keyPath, "key", "configs/key.pem", "")
	flag.Parse()

	mux := http.NewServeMux()

	mux.Handle("GET /", handlers.Middleware(handlers.NewWebHttp2Handler(web.Assets)))
	mux.Handle("POST /api/full-duplex", handlers.Middleware(handlers.NewFullDuplexHandler(flushResp, http1FullDuplex)))

	tlsCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		slog.Error("Unable to load TLS cert", slog.Any("err", err))
		exitCode = 1
		return
	}
	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	srv := http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   mux,
		TLSConfig: &tlsCfg,
	}
	if isHttp1 {
		// If TLSNextProto is not nil, HTTP/2 support is not enabled automatically
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}

	go func() {
		if err := srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Unable to listen", slog.String("add", srv.Addr), slog.Any("err", err))
			stop()
			exitCode = 1
		}
	}()

	slog.Info("Server started", slog.String("add", srv.Addr))
	<-ctx.Done()

	err = srv.Close()
	slog.Info("Server stopped")
	if err != nil {
		slog.Error("error while closing the server", slog.Any("err", err))
		exitCode = 1
	}
}
