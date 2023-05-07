// go:build !ignoreWeaverGen

package cfp

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"context"
	"dreamkast-weaver/internal/cfp/domain"
	"dreamkast-weaver/internal/cfp/value"
	"errors"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:     "dreamkast-weaver/internal/cfp/Service",
		Iface:    reflect.TypeOf((*Service)(nil)).Elem(),
		Impl:     reflect.TypeOf(ServiceImpl{}),
		ConfigFn: func(i any) any { return i.(*ServiceImpl).WithConfig.Config() },
		LocalStubFn: func(impl any, tracer trace.Tracer) any {
			return service_local_stub{impl: impl.(Service), tracer: tracer}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return service_client_stub{stub: stub, voteMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/cfp/Service", Method: "Vote"}), voteCountsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/cfp/Service", Method: "VoteCounts"})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return service_server_stub{impl: impl.(Service), addLoad: addLoad}
		},
	})
}

// Local stub implementations.

type service_local_stub struct {
	impl   Service
	tracer trace.Tracer
}

func (s service_local_stub) Vote(ctx context.Context, a0 VoteRequest) (err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "cfp.Service.Vote", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Vote(ctx, a0)
}

func (s service_local_stub) VoteCounts(ctx context.Context, a0 value.ConfName) (r0 []*domain.VoteCount, err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "cfp.Service.VoteCounts", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.VoteCounts(ctx, a0)
}

// Client stub implementations.

type service_client_stub struct {
	stub              codegen.Stub
	voteMetrics       *codegen.MethodMetrics
	voteCountsMetrics *codegen.MethodMetrics
}

func (s service_client_stub) Vote(ctx context.Context, a0 VoteRequest) (err error) {
	// Update metrics.
	start := time.Now()
	s.voteMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "cfp.Service.Vote", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.voteMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.voteMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	s.voteMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}
	s.voteMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

func (s service_client_stub) VoteCounts(ctx context.Context, a0 value.ConfName) (r0 []*domain.VoteCount, err error) {
	// Update metrics.
	start := time.Now()
	s.voteCountsMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "cfp.Service.VoteCounts", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.voteCountsMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.voteCountsMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	s.voteCountsMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 1, enc.Data(), shardKey)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}
	s.voteCountsMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = serviceweaver_dec_slice_ptr_VoteCount_dc343816(dec)
	err = dec.Error()
	return
}

// Server stub implementations.

type service_server_stub struct {
	impl    Service
	addLoad func(key uint64, load float64)
}

// GetStubFn implements the stub.Server interface.
func (s service_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "Vote":
		return s.vote
	case "VoteCounts":
		return s.voteCounts
	default:
		return nil
	}
}

func (s service_server_stub) vote(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 VoteRequest
	(&a0).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.Vote(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) voteCounts(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 value.ConfName
	(&a0).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.VoteCounts(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_slice_ptr_VoteCount_dc343816(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

// AutoMarshal implementations.

var _ codegen.AutoMarshal = &VoteRequest{}

func (x *VoteRequest) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("VoteRequest.WeaverMarshal: nil receiver"))
	}
	(x.ConfName).WeaverMarshal(enc)
	(x.TalkID).WeaverMarshal(enc)
	serviceweaver_enc_slice_byte_87461245(enc, ([]byte)(x.ClientIp))
}

func (x *VoteRequest) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("VoteRequest.WeaverUnmarshal: nil receiver"))
	}
	(&x.ConfName).WeaverUnmarshal(dec)
	(&x.TalkID).WeaverUnmarshal(dec)
	*(*[]byte)(&x.ClientIp) = serviceweaver_dec_slice_byte_87461245(dec)
}

func serviceweaver_enc_slice_byte_87461245(enc *codegen.Encoder, arg []byte) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		enc.Byte(arg[i])
	}
}

func serviceweaver_dec_slice_byte_87461245(dec *codegen.Decoder) []byte {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = dec.Byte()
	}
	return res
}

// Encoding/decoding implementations.

func serviceweaver_enc_ptr_VoteCount_316057b4(enc *codegen.Encoder, arg *domain.VoteCount) {
	if arg == nil {
		enc.Bool(false)
	} else {
		enc.Bool(true)
		(*arg).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_ptr_VoteCount_316057b4(dec *codegen.Decoder) *domain.VoteCount {
	if !dec.Bool() {
		return nil
	}
	var res domain.VoteCount
	(&res).WeaverUnmarshal(dec)
	return &res
}

func serviceweaver_enc_slice_ptr_VoteCount_dc343816(enc *codegen.Encoder, arg []*domain.VoteCount) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		serviceweaver_enc_ptr_VoteCount_316057b4(enc, arg[i])
	}
}

func serviceweaver_dec_slice_ptr_VoteCount_dc343816(dec *codegen.Decoder) []*domain.VoteCount {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]*domain.VoteCount, n)
	for i := 0; i < n; i++ {
		res[i] = serviceweaver_dec_ptr_VoteCount_316057b4(dec)
	}
	return res
}
