package graph

import (
	"dreamkast-weaver/internal/application"
	"dreamkast-weaver/internal/pkg/sqlhelper"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	stampRallyApp application.StampRallyApp
	cfpApp        application.CfpApp
	vcManager     application.ViewerCountManager
}

func New(
	sh *sqlhelper.SqlHelper,
	stampRallyApp application.StampRallyApp,
	cfpApp application.CfpApp,
	vcManager application.ViewerCountManager,
) Config {
	return Config{
		Resolvers: &Resolver{
			stampRallyApp: stampRallyApp,
			cfpApp:        cfpApp,
			vcManager:     vcManager,
		},
	}
}
