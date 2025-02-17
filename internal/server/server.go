package server

import (
	"context"
	"dreamkast-weaver/internal/graph"
	"dreamkast-weaver/internal/sqlhelper"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	gm "dreamkast-weaver/internal/graph/middleware"
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
	router := chi.NewRouter()
	router.Use(gm.ClientIP)
	router.Use(cors.Handler(corsOpts))

	sh, err := sqlhelper.NewSqlHelper(sqlhelper.NewOptionFromEnv())
	if err != nil {
		return err
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.New(sh)))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	return http.ListenAndServe(":8080", router)

}
