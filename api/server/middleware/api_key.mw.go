package middleware

import (
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"net/http"
)

// APIKeyMW ...
func APIKeyMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		logger.Infof("API KEY = %v", apiKey)
		next.ServeHTTP(w, r)
	})
}
