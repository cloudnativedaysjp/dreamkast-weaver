package main

import (
	"log"
	"os"

	"github.com/getsentry/sentry-go"
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
	envName := os.Getenv("DREAMKAST_NAMESPACE")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://c02b9c314a714cc2af37d54f0964c5ff@sentry.cloudnativedays.jp/5",
		Environment: envName,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 0.5,
		BeforeSendTransaction: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Only staging and production send error logs to sentry
			if envName == "dreamkast" || envName == "dreamkast-staging" {
				// Send the event only if it's in an appropreate namespace
				return event
			} else {
				// Don't send the transaction to Sentry otherwise
				return nil
			}
		},
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
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
