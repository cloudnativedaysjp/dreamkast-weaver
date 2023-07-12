package graph

import (
	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate
type Resolver struct {
	CfpService  cfp.Service
	DkUiService dkui.Service
}

func NewResolver(cs cfp.Service, ds dkui.Service) *Resolver {
	return &Resolver{
		CfpService:  cs,
		DkUiService: ds,
	}
}
