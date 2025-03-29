package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.New(sh, stampRallyApp, cfpApp, vcManager)))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting server...")
	return http.ListenAndServe(fmt.Sprintf(":%s", Port), router)

}
