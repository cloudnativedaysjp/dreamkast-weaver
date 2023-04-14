package graph

import (
	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/sqlhelper"
	"log"

	"github.com/ServiceWeaver/weaver"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	sqlHelper *sqlhelper.SqlHelper

	// usecases
	CfpVoter cfp.Voter
}

func NewResolver(sh *sqlhelper.SqlHelper, root weaver.Instance) *Resolver {
	voter, err := weaver.Get[cfp.Voter](root)
	if err != nil {
		log.Fatal(err)
	}

	return &Resolver{
		sqlHelper: sh,
		CfpVoter:  voter,
	}
}
