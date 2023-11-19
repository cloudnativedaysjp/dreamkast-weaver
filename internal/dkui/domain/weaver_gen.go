// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package domain

import (
	"dreamkast-weaver/internal/dkui/value"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"time"
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

// weaver.Instance checks.

// weaver.Router checks.

// Local stub implementations.

// Client stub implementations.

// Server stub implementations.

// AutoMarshal implementations.

var _ codegen.AutoMarshal = (*StampChallenge)(nil)

type __is_StampChallenge[T ~struct {
	weaver.AutoMarshal
	SlotID    value.SlotID
	Condition value.StampCondition
	UpdatedAt time.Time
}] struct{}

var _ __is_StampChallenge[StampChallenge]

func (x *StampChallenge) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("StampChallenge.WeaverMarshal: nil receiver"))
	}
	(x.SlotID).WeaverMarshal(enc)
	(x.Condition).WeaverMarshal(enc)
	enc.EncodeBinaryMarshaler(&x.UpdatedAt)
}

func (x *StampChallenge) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("StampChallenge.WeaverUnmarshal: nil receiver"))
	}
	(&x.SlotID).WeaverUnmarshal(dec)
	(&x.Condition).WeaverUnmarshal(dec)
	dec.DecodeBinaryUnmarshaler(&x.UpdatedAt)
}

var _ codegen.AutoMarshal = (*StampChallenges)(nil)

type __is_StampChallenges[T ~struct {
	weaver.AutoMarshal
	Items []StampChallenge
}] struct{}

var _ __is_StampChallenges[StampChallenges]

func (x *StampChallenges) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("StampChallenges.WeaverMarshal: nil receiver"))
	}
	serviceweaver_enc_slice_StampChallenge_a26e451e(enc, x.Items)
}

func (x *StampChallenges) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("StampChallenges.WeaverUnmarshal: nil receiver"))
	}
	x.Items = serviceweaver_dec_slice_StampChallenge_a26e451e(dec)
}

func serviceweaver_enc_slice_StampChallenge_a26e451e(enc *codegen.Encoder, arg []StampChallenge) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_StampChallenge_a26e451e(dec *codegen.Decoder) []StampChallenge {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]StampChallenge, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}

var _ codegen.AutoMarshal = (*TrackViewer)(nil)

type __is_TrackViewer[T ~struct {
	weaver.AutoMarshal
	CreatedAt time.Time
	TrackName value.TrackName
	ProfileID value.ProfileID
}] struct{}

var _ __is_TrackViewer[TrackViewer]

func (x *TrackViewer) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TrackViewer.WeaverMarshal: nil receiver"))
	}
	enc.EncodeBinaryMarshaler(&x.CreatedAt)
	(x.TrackName).WeaverMarshal(enc)
	(x.ProfileID).WeaverMarshal(enc)
}

func (x *TrackViewer) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TrackViewer.WeaverUnmarshal: nil receiver"))
	}
	dec.DecodeBinaryUnmarshaler(&x.CreatedAt)
	(&x.TrackName).WeaverUnmarshal(dec)
	(&x.ProfileID).WeaverUnmarshal(dec)
}

var _ codegen.AutoMarshal = (*TrackViewers)(nil)

type __is_TrackViewers[T ~struct {
	weaver.AutoMarshal
	Items []TrackViewer
}] struct{}

var _ __is_TrackViewers[TrackViewers]

func (x *TrackViewers) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TrackViewers.WeaverMarshal: nil receiver"))
	}
	serviceweaver_enc_slice_TrackViewer_090c5322(enc, x.Items)
}

func (x *TrackViewers) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TrackViewers.WeaverUnmarshal: nil receiver"))
	}
	x.Items = serviceweaver_dec_slice_TrackViewer_090c5322(dec)
}

func serviceweaver_enc_slice_TrackViewer_090c5322(enc *codegen.Encoder, arg []TrackViewer) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_TrackViewer_090c5322(dec *codegen.Decoder) []TrackViewer {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]TrackViewer, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}

var _ codegen.AutoMarshal = (*ViewEvent)(nil)

type __is_ViewEvent[T ~struct {
	weaver.AutoMarshal
	TrackID        value.TrackID
	TalkID         value.TalkID
	SlotID         value.SlotID
	ViewingSeconds value.ViewingSeconds
	CreatedAt      time.Time
}] struct{}

var _ __is_ViewEvent[ViewEvent]

func (x *ViewEvent) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ViewEvent.WeaverMarshal: nil receiver"))
	}
	(x.TrackID).WeaverMarshal(enc)
	(x.TalkID).WeaverMarshal(enc)
	(x.SlotID).WeaverMarshal(enc)
	(x.ViewingSeconds).WeaverMarshal(enc)
	enc.EncodeBinaryMarshaler(&x.CreatedAt)
}

func (x *ViewEvent) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ViewEvent.WeaverUnmarshal: nil receiver"))
	}
	(&x.TrackID).WeaverUnmarshal(dec)
	(&x.TalkID).WeaverUnmarshal(dec)
	(&x.SlotID).WeaverUnmarshal(dec)
	(&x.ViewingSeconds).WeaverUnmarshal(dec)
	dec.DecodeBinaryUnmarshaler(&x.CreatedAt)
}

var _ codegen.AutoMarshal = (*ViewEvents)(nil)

type __is_ViewEvents[T ~struct {
	weaver.AutoMarshal
	Items []ViewEvent
}] struct{}

var _ __is_ViewEvents[ViewEvents]

func (x *ViewEvents) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ViewEvents.WeaverMarshal: nil receiver"))
	}
	serviceweaver_enc_slice_ViewEvent_de49073f(enc, x.Items)
}

func (x *ViewEvents) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ViewEvents.WeaverUnmarshal: nil receiver"))
	}
	x.Items = serviceweaver_dec_slice_ViewEvent_de49073f(dec)
}

func serviceweaver_enc_slice_ViewEvent_de49073f(enc *codegen.Encoder, arg []ViewEvent) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_ViewEvent_de49073f(dec *codegen.Decoder) []ViewEvent {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]ViewEvent, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}

var _ codegen.AutoMarshal = (*ViewerCount)(nil)

type __is_ViewerCount[T ~struct {
	weaver.AutoMarshal
	TrackName value.TrackName
	Count     int
}] struct{}

var _ __is_ViewerCount[ViewerCount]

func (x *ViewerCount) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ViewerCount.WeaverMarshal: nil receiver"))
	}
	(x.TrackName).WeaverMarshal(enc)
	enc.Int(x.Count)
}

func (x *ViewerCount) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ViewerCount.WeaverUnmarshal: nil receiver"))
	}
	(&x.TrackName).WeaverUnmarshal(dec)
	x.Count = dec.Int()
}

var _ codegen.AutoMarshal = (*ViewerCounts)(nil)

type __is_ViewerCounts[T ~struct {
	weaver.AutoMarshal
	Items []ViewerCount
}] struct{}

var _ __is_ViewerCounts[ViewerCounts]

func (x *ViewerCounts) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ViewerCounts.WeaverMarshal: nil receiver"))
	}
	serviceweaver_enc_slice_ViewerCount_da23b7f2(enc, x.Items)
}

func (x *ViewerCounts) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ViewerCounts.WeaverUnmarshal: nil receiver"))
	}
	x.Items = serviceweaver_dec_slice_ViewerCount_da23b7f2(dec)
}

func serviceweaver_enc_slice_ViewerCount_da23b7f2(enc *codegen.Encoder, arg []ViewerCount) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_ViewerCount_da23b7f2(dec *codegen.Decoder) []ViewerCount {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]ViewerCount, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}
