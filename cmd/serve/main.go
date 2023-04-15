package main

import (
	"context"
	"dreamkast-weaver/internal/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ServiceWeaver/weaver"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize the Service Weaver application.
	root := weaver.Init(context.Background())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(root),
	}))

	opts := weaver.ListenerOptions{LocalAddress: ":" + port}
	lis, err := root.Listener("hello", opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listener available on %v\n", lis)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.Serve(lis, nil))
}
