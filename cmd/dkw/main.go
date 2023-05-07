package main

import (
	"os"

	"github.com/spf13/cobra"

	"dreamkast-weaver/cmd/dkw/dbmigrate"
	"dreamkast-weaver/cmd/dkw/serve"
)

// dkwCmd represents the base command when called without any subcommands.
var dkwCmd = &cobra.Command{
	Use:   "dkw",
	Short: "Run service and set up database",
	Long:  "Run service and set up database",
}

func init() {
	// serve
	dkwCmd.AddCommand(serve.Cmd)
	serve.Cmd.Flags().StringVarP(&serve.Port, "port", "p", "8080", "listen port")

	// dbmigrate
	dkwCmd.AddCommand(dbmigrate.Cmd)
}

func main() {
	if err := dkwCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
