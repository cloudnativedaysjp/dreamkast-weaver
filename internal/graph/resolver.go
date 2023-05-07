package graph

import (
	"github.com/ServiceWeaver/weaver"

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/dkui"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	weaver.Implements[weaver.Main]
	CfpService  weaver.Ref[cfp.Service]
	DkUiService weaver.Ref[dkui.Service]
}
