package middleware

import (
	error2 "git.teko.vn/loyalty-system/loyalty-file-processing/api/server/common/error"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"git.teko.vn/loyalty-system/loyalty-file-processing/providers/utils"
	"github.com/go-chi/render"
	"net/http"
)

// APIKeyMW ...
func APIKeyMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		logger.Infof("API KEY = %v", utils.HiddenString(apiKey, 5))
		if len(apiKey) == 0 {
			_ = render.Render(w, r, error2.ErrRenderNoPermissionRequest("Missing API Key"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
