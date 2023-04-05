package cfpsvc

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type Cfp struct {
	weaver.AutoMarshal
}

// CfpSvc component.
type T interface {
	Vote(context.Context, string) (string, error)
}

// Implementation of the CfpSvc component.
type impl struct {
	weaver.Implements[T]
}

func (s *impl) Vote(_ context.Context, str string) (string, error) {
	runes := []rune(str)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes), nil
}
