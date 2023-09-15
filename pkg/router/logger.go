package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerHTTPMiddlewareDefault(ignorePaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// skip ignore paths
			for _, prefix := range ignorePaths {
				if strings.HasPrefix(r.RequestURI, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}
			start := time.Now()

			// See if we're already using a wrapped response writer and if not make one.
			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			next.ServeHTTP(ww, r)

			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
				zap.String("path", r.RequestURI),
				zap.String("method", r.Method),
			}
			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				fields = append(fields, zap.String("request-id", reqID))
			}
			zap.L().Info("HTTP Request", fields...)
		})
	}
}
