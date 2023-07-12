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

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"
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
		if err := weaver.Run(context.Background(), serve); err != nil {
			log.Fatal(err)
		}
	},
}

type server struct {
	weaver.Implements[weaver.Main]
	CfpService  weaver.Ref[cfp.Service]
	DkUiService weaver.Ref[dkui.Service]
	//listener    weaver.Listener
}

func serve(ctx context.Context, s *server) error {
	router := chi.NewRouter()
	router.Use(gm.ClientIP)
	router.Use(cors.Handler(corsOpts))

	graphSrv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(
			s.CfpService.Get(),
			s.DkUiService.Get(),
		),
	}))

	// We do not use weaver's listener
	// because we want to inject dependencies using environment variables
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", Port))
	if err != nil {
		log.Fatal(err)
	}
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", graphSrv)

	s.Logger().Debug("Listener available", "address", lis.Addr())
	srv := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           router,
	}
	return srv.Serve(lis)
}
