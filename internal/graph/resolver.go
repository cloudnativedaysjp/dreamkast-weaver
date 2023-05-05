package graph

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ServiceWeaver/weaver"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"
	gm "dreamkast-weaver/internal/graph/middleware"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

var (
	port = "8080"
	pMu  sync.Mutex

	corsOpts = cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type", "Accept", "Authorization",
			"X-Amz-Date", "X-Api-Key", "X-Amz-Security-Token", "X-Amz-User-Agent",
		},
	}
)

func SetPort(addr string) {
	pMu.Lock()
	defer pMu.Unlock()
	port = addr
}

type Resolver struct {
	weaver.Implements[weaver.Main]
	CfpService  weaver.Ref[cfp.Service]
	DkUiService weaver.Ref[dkui.Service]
}

func Serve(ctx context.Context, r *Resolver) error {
	router := chi.NewRouter()
	router.Use(gm.ClientIP)
	router.Use(cors.Handler(corsOpts))

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: r,
	}))

	opts := weaver.ListenerOptions{LocalAddress: ":" + port}
	lis, err := r.Listener("dkw-serve", opts)
	if err != nil {
		return err
	}
	log.Printf("Listener available on %v\n", lis)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	s := http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}
	return s.Serve(lis)
}
