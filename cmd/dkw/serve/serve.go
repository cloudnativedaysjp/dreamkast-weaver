package serve

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"dreamkast-weaver/internal/server"
)

// serveCmd represents the serve command.
var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Run service",
	Long:  "Run service",
	Run: func(_ *cobra.Command, _ []string) {
		if err := server.Run(context.Background()); err != nil {
			log.Fatal(err)
		}
	},
}
