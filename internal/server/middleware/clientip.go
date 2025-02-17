package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

type ipCtxKey struct{}

func ClientIP(next http.Handler) http.Handler {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove port
		ip := strings.Split(r.RemoteAddr, ":")[0]
		ctx := context.WithValue(r.Context(), ipCtxKey{}, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	return middleware.RealIP(f)
}

func ClientIPFromContext(ctx context.Context) string {
	val := ctx.Value(ipCtxKey{})
	if val == nil {
		return ""
	}
	return val.(string)
}
