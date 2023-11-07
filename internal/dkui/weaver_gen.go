// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package dkui

import (
	"context"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"
	"errors"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

var _ codegen.LatestVersion = codegen.Version[[0][17]struct{}](`

ERROR: You generated this file with 'weaver generate' v0.17.0 (codegen
version v0.17.0). The generated code is incompatible with the version of the
github.com/ServiceWeaver/weaver module that you're using. The weaver module
version can be found in your go.mod file or by running the following command.

    go list -m github.com/ServiceWeaver/weaver

We recommend updating the weaver module and the 'weaver generate' command by
running the following.

    go get github.com/ServiceWeaver/weaver@latest
    go install github.com/ServiceWeaver/weaver/cmd/weaver@latest

Then, re-run 'weaver generate' and re-build your code. If the problem persists,
please file an issue at https://github.com/ServiceWeaver/weaver/issues.

`)

func init() {
	codegen.Register(codegen.Registration{
		Name:  "dreamkast-weaver/internal/dkui/Service",
		Iface: reflect.TypeOf((*Service)(nil)).Elem(),
		Impl:  reflect.TypeOf(ServiceImpl{}),
		LocalStubFn: func(impl any, caller string, tracer trace.Tracer) any {
			return service_local_stub{impl: impl.(Service), tracer: tracer, createViewEventMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "CreateViewEvent", Remote: false}), listViewerCountsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "ListViewerCounts", Remote: false}), stampChallengesMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "StampChallenges", Remote: false}), stampOnSiteMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "StampOnSite", Remote: false}), stampOnlineMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "StampOnline", Remote: false}), viewingEventsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "ViewingEvents", Remote: false}), viewingTrackMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "ViewingTrack", Remote: false})}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return service_client_stub{stub: stub, createViewEventMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "CreateViewEvent", Remote: true}), listViewerCountsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "ListViewerCounts", Remote: true}), stampChallengesMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "StampChallenges", Remote: true}), stampOnSiteMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "StampOnSite", Remote: true}), stampOnlineMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "StampOnline", Remote: true}), viewingEventsMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "ViewingEvents", Remote: true}), viewingTrackMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "dreamkast-weaver/internal/dkui/Service", Method: "ViewingTrack", Remote: true})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return service_server_stub{impl: impl.(Service), addLoad: addLoad}
		},
		RefData: "",
	})
}

// weaver.Instance checks.
var _ weaver.InstanceOf[Service] = (*ServiceImpl)(nil)

// weaver.Router checks.
var _ weaver.Unrouted = (*ServiceImpl)(nil)

// Local stub implementations.

type service_local_stub struct {
	impl                    Service
	tracer                  trace.Tracer
	createViewEventMetrics  *codegen.MethodMetrics
	listViewerCountsMetrics *codegen.MethodMetrics
	stampChallengesMetrics  *codegen.MethodMetrics
	stampOnSiteMetrics      *codegen.MethodMetrics
	stampOnlineMetrics      *codegen.MethodMetrics
	viewingEventsMetrics    *codegen.MethodMetrics
	viewingTrackMetrics     *codegen.MethodMetrics
}

// Check that service_local_stub implements the Service interface.
var _ Service = (*service_local_stub)(nil)

func (s service_local_stub) CreateViewEvent(ctx context.Context, a0 Profile, a1 CreateViewEventRequest) (err error) {
	// Update metrics.
	begin := s.createViewEventMetrics.Begin()
	defer func() { s.createViewEventMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "dkui.Service.CreateViewEvent", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.CreateViewEvent(ctx, a0, a1)
}

func (s service_local_stub) ListViewerCounts(ctx context.Context) (r0 *domain.ViewerCounts, err error) {
	// Update metrics.
	begin := s.listViewerCountsMetrics.Begin()
	defer func() { s.listViewerCountsMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "dkui.Service.ListViewerCounts", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.ListViewerCounts(ctx)
}

func (s service_local_stub) StampChallenges(ctx context.Context, a0 Profile) (r0 *domain.StampChallenges, err error) {
	// Update metrics.
	begin := s.stampChallengesMetrics.Begin()
	defer func() { s.stampChallengesMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "dkui.Service.StampChallenges", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.StampChallenges(ctx, a0)
}

func (s service_local_stub) StampOnSite(ctx context.Context, a0 Profile, a1 StampRequest) (err error) {
	// Update metrics.
	begin := s.stampOnSiteMetrics.Begin()
	defer func() { s.stampOnSiteMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "dkui.Service.StampOnSite", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.StampOnSite(ctx, a0, a1)
}

func (s service_local_stub) StampOnline(ctx context.Context, a0 Profile, a1 value.SlotID) (err error) {
	// Update metrics.
	begin := s.stampOnlineMetrics.Begin()
	defer func() { s.stampOnlineMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "dkui.Service.StampOnline", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.StampOnline(ctx, a0, a1)
}

func (s service_local_stub) ViewingEvents(ctx context.Context, a0 Profile) (r0 *domain.ViewEvents, err error) {
	// Update metrics.
	begin := s.viewingEventsMetrics.Begin()
	defer func() { s.viewingEventsMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "dkui.Service.ViewingEvents", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.ViewingEvents(ctx, a0)
}

func (s service_local_stub) ViewingTrack(ctx context.Context, a0 value.ProfileID, a1 value.TrackName) (err error) {
	// Update metrics.
	begin := s.viewingTrackMetrics.Begin()
	defer func() { s.viewingTrackMetrics.End(begin, err != nil, 0, 0) }()
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "dkui.Service.ViewingTrack", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.ViewingTrack(ctx, a0, a1)
}

// Client stub implementations.

type service_client_stub struct {
	stub                    codegen.Stub
	createViewEventMetrics  *codegen.MethodMetrics
	listViewerCountsMetrics *codegen.MethodMetrics
	stampChallengesMetrics  *codegen.MethodMetrics
	stampOnSiteMetrics      *codegen.MethodMetrics
	stampOnlineMetrics      *codegen.MethodMetrics
	viewingEventsMetrics    *codegen.MethodMetrics
	viewingTrackMetrics     *codegen.MethodMetrics
}

// Check that service_client_stub implements the Service interface.
var _ Service = (*service_client_stub)(nil)

func (s service_client_stub) CreateViewEvent(ctx context.Context, a0 Profile, a1 CreateViewEventRequest) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.createViewEventMetrics.Begin()
	defer func() { s.createViewEventMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "dkui.Service.CreateViewEvent", trace.WithSpanKind(trace.SpanKindClient))
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
		}
		span.End()

	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	(a1).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

func (s service_client_stub) ListViewerCounts(ctx context.Context) (r0 *domain.ViewerCounts, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.listViewerCountsMetrics.Begin()
	defer func() { s.listViewerCountsMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "dkui.Service.ListViewerCounts", trace.WithSpanKind(trace.SpanKindClient))
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
		}
		span.End()

	}()

	var shardKey uint64

	// Call the remote method.
	var results []byte
	results, err = s.stub.Run(ctx, 1, nil, shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = serviceweaver_dec_ptr_ViewerCounts_da8696ef(dec)
	err = dec.Error()
	return
}

func (s service_client_stub) StampChallenges(ctx context.Context, a0 Profile) (r0 *domain.StampChallenges, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.stampChallengesMetrics.Begin()
	defer func() { s.stampChallengesMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "dkui.Service.StampChallenges", trace.WithSpanKind(trace.SpanKindClient))
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
		}
		span.End()

	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 2, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = serviceweaver_dec_ptr_StampChallenges_fd60f90f(dec)
	err = dec.Error()
	return
}

func (s service_client_stub) StampOnSite(ctx context.Context, a0 Profile, a1 StampRequest) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.stampOnSiteMetrics.Begin()
	defer func() { s.stampOnSiteMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "dkui.Service.StampOnSite", trace.WithSpanKind(trace.SpanKindClient))
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
		}
		span.End()

	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	(a1).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 3, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

func (s service_client_stub) StampOnline(ctx context.Context, a0 Profile, a1 value.SlotID) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.stampOnlineMetrics.Begin()
	defer func() { s.stampOnlineMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "dkui.Service.StampOnline", trace.WithSpanKind(trace.SpanKindClient))
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
		}
		span.End()

	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	(a1).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 4, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

func (s service_client_stub) ViewingEvents(ctx context.Context, a0 Profile) (r0 *domain.ViewEvents, err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.viewingEventsMetrics.Begin()
	defer func() { s.viewingEventsMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "dkui.Service.ViewingEvents", trace.WithSpanKind(trace.SpanKindClient))
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
		}
		span.End()

	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 5, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = serviceweaver_dec_ptr_ViewEvents_fd9e24a3(dec)
	err = dec.Error()
	return
}

func (s service_client_stub) ViewingTrack(ctx context.Context, a0 value.ProfileID, a1 value.TrackName) (err error) {
	// Update metrics.
	var requestBytes, replyBytes int
	begin := s.viewingTrackMetrics.Begin()
	defer func() { s.viewingTrackMetrics.End(begin, err != nil, requestBytes, replyBytes) }()

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "dkui.Service.ViewingTrack", trace.WithSpanKind(trace.SpanKindClient))
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
		}
		span.End()

	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	(a1).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	requestBytes = len(enc.Data())
	var results []byte
	results, err = s.stub.Run(ctx, 6, enc.Data(), shardKey)
	replyBytes = len(results)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}

	// Decode the results.
	dec := codegen.NewDecoder(results)
	err = dec.Error()
	return
}

// Server stub implementations.

type service_server_stub struct {
	impl    Service
	addLoad func(key uint64, load float64)
}

// Check that service_server_stub implements the codegen.Server interface.
var _ codegen.Server = (*service_server_stub)(nil)

// GetStubFn implements the codegen.Server interface.
func (s service_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "CreateViewEvent":
		return s.createViewEvent
	case "ListViewerCounts":
		return s.listViewerCounts
	case "StampChallenges":
		return s.stampChallenges
	case "StampOnSite":
		return s.stampOnSite
	case "StampOnline":
		return s.stampOnline
	case "ViewingEvents":
		return s.viewingEvents
	case "ViewingTrack":
		return s.viewingTrack
	default:
		return nil
	}
}

func (s service_server_stub) createViewEvent(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 Profile
	(&a0).WeaverUnmarshal(dec)
	var a1 CreateViewEventRequest
	(&a1).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.CreateViewEvent(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) listViewerCounts(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.ListViewerCounts(ctx)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_ptr_ViewerCounts_da8696ef(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) stampChallenges(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 Profile
	(&a0).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.StampChallenges(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_ptr_StampChallenges_fd60f90f(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) stampOnSite(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 Profile
	(&a0).WeaverUnmarshal(dec)
	var a1 StampRequest
	(&a1).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.StampOnSite(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) stampOnline(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 Profile
	(&a0).WeaverUnmarshal(dec)
	var a1 value.SlotID
	(&a1).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.StampOnline(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) viewingEvents(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 Profile
	(&a0).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.ViewingEvents(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	serviceweaver_enc_ptr_ViewEvents_fd9e24a3(enc, r0)
	enc.Error(appErr)
	return enc.Data(), nil
}

func (s service_server_stub) viewingTrack(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 value.ProfileID
	(&a0).WeaverUnmarshal(dec)
	var a1 value.TrackName
	(&a1).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	appErr := s.impl.ViewingTrack(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.Error(appErr)
	return enc.Data(), nil
}

// AutoMarshal implementations.

var _ codegen.AutoMarshal = (*CreateViewEventRequest)(nil)

type __is_CreateViewEventRequest[T ~struct {
	weaver.AutoMarshal
	TrackID value.TrackID
	TalkID  value.TalkID
	SlotID  value.SlotID
}] struct{}

var _ __is_CreateViewEventRequest[CreateViewEventRequest]

func (x *CreateViewEventRequest) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("CreateViewEventRequest.WeaverMarshal: nil receiver"))
	}
	(x.TrackID).WeaverMarshal(enc)
	(x.TalkID).WeaverMarshal(enc)
	(x.SlotID).WeaverMarshal(enc)
}

func (x *CreateViewEventRequest) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("CreateViewEventRequest.WeaverUnmarshal: nil receiver"))
	}
	(&x.TrackID).WeaverUnmarshal(dec)
	(&x.TalkID).WeaverUnmarshal(dec)
	(&x.SlotID).WeaverUnmarshal(dec)
}

var _ codegen.AutoMarshal = (*Profile)(nil)

type __is_Profile[T ~struct {
	weaver.AutoMarshal
	ID       value.ProfileID
	ConfName value.ConfName
}] struct{}

var _ __is_Profile[Profile]

func (x *Profile) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Profile.WeaverMarshal: nil receiver"))
	}
	(x.ID).WeaverMarshal(enc)
	(x.ConfName).WeaverMarshal(enc)
}

func (x *Profile) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Profile.WeaverUnmarshal: nil receiver"))
	}
	(&x.ID).WeaverUnmarshal(dec)
	(&x.ConfName).WeaverUnmarshal(dec)
}

var _ codegen.AutoMarshal = (*StampRequest)(nil)

type __is_StampRequest[T ~struct {
	weaver.AutoMarshal
	TrackID value.TrackID
	TalkID  value.TalkID
	SlotID  value.SlotID
}] struct{}

var _ __is_StampRequest[StampRequest]

func (x *StampRequest) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("StampRequest.WeaverMarshal: nil receiver"))
	}
	(x.TrackID).WeaverMarshal(enc)
	(x.TalkID).WeaverMarshal(enc)
	(x.SlotID).WeaverMarshal(enc)
}

func (x *StampRequest) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("StampRequest.WeaverUnmarshal: nil receiver"))
	}
	(&x.TrackID).WeaverUnmarshal(dec)
	(&x.TalkID).WeaverUnmarshal(dec)
	(&x.SlotID).WeaverUnmarshal(dec)
}

// Encoding/decoding implementations.

func serviceweaver_enc_ptr_ViewerCounts_da8696ef(enc *codegen.Encoder, arg *domain.ViewerCounts) {
	if arg == nil {
		enc.Bool(false)
	} else {
		enc.Bool(true)
		(*arg).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_ptr_ViewerCounts_da8696ef(dec *codegen.Decoder) *domain.ViewerCounts {
	if !dec.Bool() {
		return nil
	}
	var res domain.ViewerCounts
	(&res).WeaverUnmarshal(dec)
	return &res
}

func serviceweaver_enc_ptr_StampChallenges_fd60f90f(enc *codegen.Encoder, arg *domain.StampChallenges) {
	if arg == nil {
		enc.Bool(false)
	} else {
		enc.Bool(true)
		(*arg).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_ptr_StampChallenges_fd60f90f(dec *codegen.Decoder) *domain.StampChallenges {
	if !dec.Bool() {
		return nil
	}
	var res domain.StampChallenges
	(&res).WeaverUnmarshal(dec)
	return &res
}

func serviceweaver_enc_ptr_ViewEvents_fd9e24a3(enc *codegen.Encoder, arg *domain.ViewEvents) {
	if arg == nil {
		enc.Bool(false)
	} else {
		enc.Bool(true)
		(*arg).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_ptr_ViewEvents_fd9e24a3(dec *codegen.Decoder) *domain.ViewEvents {
	if !dec.Bool() {
		return nil
	}
	var res domain.ViewEvents
	(&res).WeaverUnmarshal(dec)
	return &res
}
