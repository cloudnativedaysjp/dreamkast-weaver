package graph

import (
	"dreamkast-weaver/internal/sqlhelper"
	"dreamkast-weaver/internal/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	dkUiSrv usecase.DkUiService
	cfpSrv  usecase.CfpService
}

func New(sh *sqlhelper.SqlHelper) Config {
	return Config{
		Resolvers: &Resolver{
			dkUiSrv: usecase.NewDkUiService(sh),
			cfpSrv:  usecase.NewCFPService(sh),
		},
	}
}
