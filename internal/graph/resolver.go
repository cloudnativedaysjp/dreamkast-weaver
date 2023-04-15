package graph

import (
	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"
	"log"

	"github.com/ServiceWeaver/weaver"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	CfpVoter    cfp.Voter
	DkUiService dkui.Service
}

func NewResolver(root weaver.Instance) *Resolver {
	cfp, err := weaver.Get[cfp.Voter](root)
	if err != nil {
		log.Fatal(err)
	}

	dkui, err := weaver.Get[dkui.Service](root)
	if err != nil {
		log.Fatal(err)
	}

	return &Resolver{
		CfpVoter:    cfp,
		DkUiService: dkui,
	}
}
