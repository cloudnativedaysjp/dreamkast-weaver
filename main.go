package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ServiceWeaver/weaver"

	"dreamkast-weaver/service/frontend"
)

//go:generate weaver generate ./...

var localAddr = flag.String("local_addr", ":12345", "Local address")

func main() {
	// Initialize the Service Weaver application.
	ctx := context.Background()
	root := weaver.Init(ctx)
	server, err := frontend.NewServer(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating frontend: ", err)
		os.Exit(1)
	}
	if err := server.Run(*localAddr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
