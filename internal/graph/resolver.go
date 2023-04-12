package graph

import (
	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/sqlhelper"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	sqlHelper *sqlhelper.SqlHelper

	// usecases
	CfpVoter cfp.Voter
}

func NewResolver(sh *sqlhelper.SqlHelper) *Resolver {
	return &Resolver{
		sqlHelper: sh,
		CfpVoter:  cfp.NewVoter(),
	}
}
