package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

//go:embed configs/cert.pem
var cert []byte

//go:embed configs/key.pem
var key []byte

//go:embed web
var web embed.FS

var _ http.Handler = pushHandler{}

type pushHandler struct {
	next http.Handler
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	var isHttp1 bool
	flag.BoolVar(&isHttp1, "http1", false, "")
	flag.Parse()

	mux := http.NewServeMux()
	content, err := fs.Sub(web, "web")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServerFS(content)

	mux.Handle("GET /", middleware(pushHandler{next: fileServer}))
	mux.Handle("POST /api/full-duplex", middleware(http.HandlerFunc(fullDuplexHandler)))

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
	if isHttp1 {
		srv.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
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

// Write implements io.Writer.
// func (f flushWriter) Write(p []byte) (int, error) {
// 	n, err := f.w.Write(p)
// 	if err != nil {
// 		return n, err
// 	}
// 	rc := http.NewResponseController(f.w)
// 	rc.EnableFullDuplex()
// 	if err := rc.Flush(); err != nil {
// 		if errors.Is(err, http.ErrNotSupported) {
// 			log.Println("Flush is not supported")
// 		} else {
// 			log.Panicf("unexpected error during flushing: %s", err)
// 		}
// 	}
// 	return n, err
// }

func (p pushHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		if pusher, ok := w.(http.Pusher); ok {
			pusher.Push("/index.css", nil)
		}
	}
	p.next.ServeHTTP(w, r)
}

func fullDuplexHandler(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	if r.ProtoMajor == 1 {
		log.Print("unsuported")
		w.WriteHeader(http.StatusHTTPVersionNotSupported)
		return
	}
	if r.ProtoMajor == 1 {
		log.Printf("trying to enable full duplex for %s", r.Proto)
		if err := rc.EnableFullDuplex(); err != nil {
			log.Printf("failed to enable full duplex: %s", err)
			if errors.Is(err, http.ErrNotSupported) {
				w.WriteHeader(http.StatusHTTPVersionNotSupported)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}

	ctx := r.Context()
	go func() {
		bufReader := bufio.NewReader(r.Body)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				req, err := bufReader.ReadString('\n')
				if err != nil {
					log.Printf("unable to ready body: %s", err)
				}
				log.Print(req)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.Write([]byte("server resp\n"))
			if err := rc.Flush(); err != nil {
				log.Printf("failed to flush: %s", err)
				if errors.Is(err, http.ErrNotSupported) {
					w.WriteHeader(http.StatusHTTPVersionNotSupported)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %T", r.Proto, r.Method, r.URL.Path, w)

		w.Header().Set("Cache-Control", "no-cache; no-store; must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)

		log.Print("Request served")
	})
}
