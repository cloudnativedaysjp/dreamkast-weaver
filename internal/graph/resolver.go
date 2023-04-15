package graph

import (
	"log"

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"
	"github.com/ServiceWeaver/weaver"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	CfpService  cfp.Service
	DkUiService dkui.Service
}

func NewResolver(root weaver.Instance) *Resolver {
	cfp, err := weaver.Get[cfp.Service](root)
	if err != nil {
		log.Fatal(err)
	}

	dkui, err := weaver.Get[dkui.Service](root)
	if err != nil {
		log.Fatal(err)
	}

	return &Resolver{
		CfpService:  cfp,
		DkUiService: dkui,
	}
}
