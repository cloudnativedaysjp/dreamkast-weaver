package serve

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ServiceWeaver/weaver"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"

	"dreamkast-weaver/internal/graph"
	gm "dreamkast-weaver/internal/graph/middleware"
)

var Port string

var (
	corsOpts = cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Content-Type", "Accept", "Authorization",
			"X-Amz-Date", "X-Api-Key", "X-Amz-Security-Token", "X-Amz-User-Agent",
		},
	}
)

// Cmd represents the serve command.
var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Run service",
	Long:  "Run service",
	Run: func(cmd *cobra.Command, args []string) {

		router := chi.NewRouter()
		router.Use(gm.ClientIP)
		router.Use(cors.Handler(corsOpts))

		// Initialize the Service Weaver application.
		root := weaver.Init(context.Background())

		srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
			Resolvers: graph.NewResolver(root),
		}))

		opts := weaver.ListenerOptions{LocalAddress: ":" + Port}
		lis, err := root.Listener("dkw-serve", opts)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listener available on %v\n", lis)

		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
		router.Handle("/query", srv)

		log.Printf("connect to http://localhost:%s/ for GraphQL playground", Port)
		s := http.Server{
			ReadHeaderTimeout: 5 * time.Second,
			Handler:           router,
		}
		log.Fatal(s.Serve(lis))
	},
}
