package middleware

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

// LoggerMW ...
func LoggerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		startReqTime := time.Now()
		defer func() {
			logger.Infof("%s %s%s%s %s %d %dB in %s",
				r.Method,
				r.URL.Scheme,
				r.Host,
				r.URL.Path,
				r.Proto,
				ww.Status(),
				ww.BytesWritten(),
				time.Since(startReqTime),
			)
		}()
		next.ServeHTTP(ww, r)
	})
}
