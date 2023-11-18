package serve

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ServiceWeaver/weaver"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"

	"dreamkast-weaver/internal/graph"
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

// serveCmd represents the serve command.
var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Run service",
	Long:  "Run service",
	Run: func(_ *cobra.Command, _ []string) {
		if err := weaver.Run(context.Background(), serve); err != nil {
			log.Fatal(err)
		}
	},
}

func serve(ctx context.Context, r *graph.Resolver) error {
	router := chi.NewRouter()
	router.Use(gm.ClientIP)
	router.Use(cors.Handler(corsOpts))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: r,
	}))

	var lis net.Listener
	if Port == "" {
		lis = r.Graphql
	} else {
		// If Port is defined, it takes precedence over weaver.toml
		l, err := net.Listen("tcp", fmt.Sprintf(":%s", Port))
		if err != nil {
			return err
		}
		lis = l
	}
	r.Logger().Debug("Listener available", "address", lis.Addr())

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	s := http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}
	return s.Serve(lis)
}
