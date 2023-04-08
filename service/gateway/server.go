package gateway

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ServiceWeaver/weaver"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"dreamkast-weaver/graph"
	"dreamkast-weaver/service/cfpsvc"
)

const defaultPort = "8080"

type Server struct {
	handler  http.Handler
	root     weaver.Instance
	hostname string

	cfpSvc cfpsvc.T
}

func NewServer(root weaver.Instance) (*Server, error) {
	// Setup the services.
	cfpSvc, err := weaver.Get[cfpsvc.T](root)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		root.Logger().Debug(`cannot get hostname for frontend: using "unknown"`)
		hostname = "unknown"
	}

	// Create the server.
	s := &Server{
		root:     root,
		hostname: hostname,
		cfpSvc:   cfpSvc,
	}

	r := http.NewServeMux()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", weaver.InstrumentHandler("query", srv))

	// Set handler and return.
	var handler http.Handler = r
	handler = otelhttp.NewHandler(handler, "http") // add tracing
	s.handler = handler

	return s, nil
}

func (s *Server) Run(localAddr string) error {
	lis, err := s.root.Listener("boutique", weaver.ListenerOptions{LocalAddress: localAddr})
	if err != nil {
		return err
	}
	s.root.Logger().Debug("Frontend available", "addr", lis)
	return http.Serve(lis, s.handler)
}
