package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"git.teko.vn/loyalty-system/loyalty-file-processing/pkg/logger"
	"net/http"
)

const userHeader = "X-USER"

// UserMW ...
func UserMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get header
		userRaw := r.Header.Get(userHeader)
		logger.Infof("USER = %v", userRaw)
		user := &User{}
		if len(userRaw) == 0 {
			msg := fmt.Sprintf("Missing %v", userHeader)
			logger.Warn(msg)
			//_ = render.Render(w, r, error2.ErrRenderInvalidRequest(errors.New(msg)))
			user = &User{
				Sub:   "unknown_sub",
				Name:  "unknown_user",
				Email: "unknown_user@teko.vn",
			}
		} else {
			// Map header to User
			user := &User{}
			err := json.Unmarshal([]byte(userRaw), &user)
			if err != nil {
				logger.Errorf("Cannot unmarshal user info in header: err=%v, rawUser=%v", err, userRaw)
				// todo return error
				user = &User{
					Sub:   "unknown_sub",
					Name:  "unknown_user",
					Email: "unknown_user@teko.vn",
				}
			}
		}

		// Set to context
		ctx := r.Context()
		ctx = SetUserToContext(ctx, *user)
		r = r.WithContext(ctx)

		// Next
		next.ServeHTTP(w, r)
	})
}

// User is the data of IAM User, that is got in BFF layer
type User struct {
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type contextKey string

var contextKeyUser = contextKey("user-attribute")

// SetUserToContext ...
func SetUserToContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, contextKeyUser, user)
}

// GetUserFromContext ... get user from context, should get in Service layer
func GetUserFromContext(ctx context.Context) User {
	user, ok := ctx.Value(contextKeyUser).(User)
	if !ok {
		return User{}
	}
	return user
}
