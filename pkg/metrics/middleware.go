package metrics

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// BuildRequestMiddleware builds middleware that produces prometheus metrics
func BuildRequestMiddleware(registry *prometheus.Registry) func(next http.Handler) http.Handler {
	httpRequestCounter := buildHTTPRequestCounterCollector()

	err := registry.Register(httpRequestCounter)
	if err != nil {
		// TODO: add custom logs
		log.Printf("BuildRequestMiddleware failed: %w", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := newResponseWriter(w)

			next.ServeHTTP(rw, r)
			httpRequestCounter.WithLabelValues(r.Method, r.RequestURI, strconv.Itoa(rw.statusCode)).Inc()
		})
	}
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriterInterceptor {
	return &responseWriterInterceptor{w, http.StatusOK}
}

func (rw *responseWriterInterceptor) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterInterceptor) Write(p []byte) (int, error) {
	return rw.ResponseWriter.Write(p)
}

func (rw *responseWriterInterceptor) Flush() {
	f, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	f.Flush()
}
