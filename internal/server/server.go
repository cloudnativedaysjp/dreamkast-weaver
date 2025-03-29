package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vektah/gqlparser/v2/ast"

	"dreamkast-weaver/internal/application"
	"dreamkast-weaver/internal/pkg/sqlhelper"
	"dreamkast-weaver/internal/server/graph"
	gm "dreamkast-weaver/internal/server/middleware"
)

var (
	Port string

	corsOpts = cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}
)

func Run(ctx context.Context) error {
	fmt.Println("Initializing server...")
	router := chi.NewRouter()
	router.Use(gm.ClientIP)
	router.Use(cors.Handler(corsOpts))

	sh, err := sqlhelper.NewSqlHelper(sqlhelper.NewOptionFromEnv())
	if err != nil {
		return err
	}

	stampRallyApp := application.NewStampRallyApp(sh)
	vcManager := application.NewViewerCountManager(sh)
	cfpApp := application.NewCfpApp(sh)

	srv := newServer(graph.NewExecutableSchema(graph.New(sh, stampRallyApp, cfpApp, vcManager)))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Starting server...")
	return server.ListenAndServe()
}

func newServer(es graphql.ExecutableSchema) *handler.Server {
	srv := handler.New(es)
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}
