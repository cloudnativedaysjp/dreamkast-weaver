package graph

import (
	"context"
	"log"
	"net/http"

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"

	"github.com/ServiceWeaver/weaver"
	"github.com/tomasen/realip"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	CfpService  cfp.Service
	DkUiService dkui.Service
}

func NewResolver(root weaver.Instance) *Resolver {
	cfp, err := weaver.Get[cfp.Service](root)
	if err != nil {
		log.Fatal(err)
	}

	dkui, err := weaver.Get[dkui.Service](root)
	if err != nil {
		log.Fatal(err)
	}

	return &Resolver{
		CfpService:  cfp,
		DkUiService: dkui,
	}
}

type ipCtxKey struct{}

func NewIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)
		ctx := context.WithValue(r.Context(), ipCtxKey{}, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClientIP(ctx context.Context) string {
	val := ctx.Value(ipCtxKey{})
	if val == nil {
		return ""
	}

	return val.(string)
}
