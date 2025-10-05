package handlers

import (
	"errors"
	"go-code-patterns/http2/pkg/respctrl"
	"io/fs"
	"log/slog"
	"net/http"
)

var _ http.Handler = webHttp2Handler{}

type webHttp2Handler struct {
	fileServer http.Handler
}

func NewWebHttp2Handler(content fs.FS) http.Handler {
	return webHttp2Handler{
		fileServer: http.FileServerFS(content),
	}
}

func (wh webHttp2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		rc := respctrl.NewImprovedResponseController(w)
		if err := rc.Push("/index.css", nil); err != nil {
			if errors.Is(err, http.ErrNotSupported) {
				slog.Warn("Pushing is not supported")
			} else {
				slog.Error("Error while pushing", slog.Any("err", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			slog.Info("Successfully pushed", slog.String("resource", "/index.css"))
		}
	}
	wh.fileServer.ServeHTTP(w, r)
}
