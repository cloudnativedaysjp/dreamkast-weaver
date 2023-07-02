package serve

import (
	"context"
	"log"
	"net/http"

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
		AllowedHeaders: []string{
			"Content-Type", "Accept", "Authorization",
			"X-Amz-Date", "X-Api-Key", "X-Amz-Security-Token", "X-Amz-User-Agent",
		},
	}
)

// serveCmd represents the serve command.
var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Run service",
	Long:  "Run service",
	Run: func(_ *cobra.Command, _ []string) {
		if err := weaver.Run(context.Background()); err != nil {
			log.Fatal(err)
		}
	},
}

type server struct {
	weaver.Implements[weaver.Main]
	resolver weaver.Ref[graph.Resolver]
}

func (s *server) Main(ctx context.Context) error {
	router := chi.NewRouter()
	router.Use(gm.ClientIP)
	router.Use(cors.Handler(corsOpts))

	r := s.resolver.Get()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &r,
	}))

	opts := weaver.ListenerOptions{LocalAddress: ":" + Port}
	lis, err := s.Listener("dkw-serve", opts)
	if err != nil {
		return err
	}
	log.Printf("Listener available on %v\n", lis)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", Port)
	return http.Serve(lis, router)
}
