package main

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

//go:embed configs/cert.pem
var cert []byte

//go:embed configs/key.pem
var key []byte

//go:embed web
var web embed.FS

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	mux := http.NewServeMux()
	content, err := fs.Sub(web, "web")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServerFS(content)
	mux.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			if pusher, ok := w.(http.Pusher); ok {
				pusher.Push("/index.css", nil)
			}
		}
		fileServer.ServeHTTP(w, r)
	})))

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}
	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}

	srv := http.Server{
		Addr:      ":8080",
		Handler:   mux,
		TLSConfig: &tlsCfg,
	}

	go func() {
		if err := srv.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			log.Print(err)
			stop()
		}
	}()
	log.Printf("Listening on %s", srv.Addr)

	<-ctx.Done()

	err = srv.Close()
	if err != nil {
		log.Print(err)
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %T", r.Proto, r.Method, r.URL.Path, w)

		w.Header().Set("Cache-Control", "no-cache; no-store; must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}
