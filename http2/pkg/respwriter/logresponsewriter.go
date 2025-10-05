package respwriter

import "net/http"

// Ensure loggedResponseWriter has `Unwrap` method
// as per doc https://go.dev/src/net/http/responsecontroller.go
var _ interface{ Unwrap() http.ResponseWriter } = (*LogResponseWriter)(nil)

type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode   int
	WroteHeader  bool
	WrittenBytes int64
}

// WriteHeader is http.ResponseWriter implementation that captures status code
func (l *LogResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	if statusCode >= 200 && !l.WroteHeader {
		l.StatusCode = statusCode
		l.WroteHeader = true
	}
}

// Write is http.ResponseWriter implementation that captures the number of written bytes
func (l *LogResponseWriter) Write(b []byte) (int, error) {
	writtenBytes, err := l.ResponseWriter.Write(b)
	if !l.WroteHeader {
		l.StatusCode = http.StatusOK
		l.WroteHeader = true
	}
	l.WrittenBytes += int64(writtenBytes)
	return writtenBytes, err
}

// Unwrap returns the original http.ResponseWriter
// as per https://go.dev/src/net/http/responsecontroller.go
func (l *LogResponseWriter) Unwrap() http.ResponseWriter {
	return l.ResponseWriter
}
