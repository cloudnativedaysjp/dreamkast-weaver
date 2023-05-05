package serve

import (
	"context"
	"log"

	"github.com/ServiceWeaver/weaver"
	"github.com/spf13/cobra"

	"dreamkast-weaver/internal/graph"
)

var Port string

// serveCmd represents the serve command.
var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Run service",
	Long:  "Run service",
	Run: func(_ *cobra.Command, _ []string) {
		graph.SetPort(Port)
		if err := weaver.Run(context.Background(), graph.Serve); err != nil {
			log.Fatal(err)
		}
	},
}
