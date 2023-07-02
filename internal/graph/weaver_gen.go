// go:build !ignoreWeaverGen

package graph

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"context"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

var _ codegen.LatestVersion = codegen.Version[[0][11]struct{}]("You used 'weaver generate' codegen version 0.11.0, but you built your code with an incompatible weaver module version. Try upgrading 'weaver generate' and re-running it.")

func init() {
	codegen.Register(codegen.Registration{
		Name:  "dreamkast-weaver/internal/graph/ResolverIF",
		Iface: reflect.TypeOf((*ResolverIF)(nil)).Elem(),
		Impl:  reflect.TypeOf(Resolver{}),
		LocalStubFn: func(impl any, tracer trace.Tracer) any {
			return resolverIF_local_stub{impl: impl.(ResolverIF), tracer: tracer}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any { return resolverIF_client_stub{stub: stub} },
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return resolverIF_server_stub{impl: impl.(ResolverIF), addLoad: addLoad}
		},
		RefData: "⟦15494570:wEaVeReDgE:dreamkast-weaver/internal/graph/ResolverIF→dreamkast-weaver/internal/cfp/Service⟧\n⟦79dd111a:wEaVeReDgE:dreamkast-weaver/internal/graph/ResolverIF→dreamkast-weaver/internal/dkui/Service⟧\n",
	})
}

// weaver.Instance checks.
var _ weaver.InstanceOf[ResolverIF] = (*Resolver)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*Resolver)(nil)

// Local stub implementations.

type resolverIF_local_stub struct {
	impl   ResolverIF
	tracer trace.Tracer
}

// Check that resolverIF_local_stub implements the ResolverIF interface.
var _ ResolverIF = (*resolverIF_local_stub)(nil)

// Client stub implementations.

type resolverIF_client_stub struct {
	stub codegen.Stub
}

// Check that resolverIF_client_stub implements the ResolverIF interface.
var _ ResolverIF = (*resolverIF_client_stub)(nil)

// Server stub implementations.

type resolverIF_server_stub struct {
	impl    ResolverIF
	addLoad func(key uint64, load float64)
}

// Check that resolverIF_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*resolverIF_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s resolverIF_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	default:
		return nil
	}
}
