package handlers

import (
	"go-code-patterns/http2/pkg/respwriter"
	"log/slog"
	"net/http"
	"reflect"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proto := r.Proto
		method := r.Method
		path := r.URL.Path

		w.Header().Set("Cache-Control", "no-cache; no-store; must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		lrw := respwriter.LogResponseWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(&lrw, r)

		slog.Info("Request served",
			slog.String("protocol", proto),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status_code", lrw.StatusCode),
			slog.Int64("response_size", lrw.WrittenBytes),
			slog.String("responsewriter_type", reflect.TypeOf(w).String()),
		)
	})
}
