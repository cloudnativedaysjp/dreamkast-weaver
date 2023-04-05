package cfpsvc

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

// CfpSvc component.
type Cfp interface {
	Vote(context.Context, string) (string, error)
}

// Implementation of the CfpSvc component.
type cfp struct {
	weaver.Implements[Cfp]
}

func (r *cfp) Vote(_ context.Context, s string) (string, error) {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes), nil
}
