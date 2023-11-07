// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package value

import (
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
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

var _ codegen.AutoMarshal = (*ConfName)(nil)

type __is_ConfName[T ~struct {
	weaver.AutoMarshal
	value ConferenceKind
}] struct{}

var _ __is_ConfName[ConfName]

func (x *ConfName) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ConfName.WeaverMarshal: nil receiver"))
	}
	enc.String((string)(x.value))
}

func (x *ConfName) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ConfName.WeaverUnmarshal: nil receiver"))
	}
	*(*string)(&x.value) = dec.String()
}

var _ codegen.AutoMarshal = (*ProfileID)(nil)

type __is_ProfileID[T ~struct {
	weaver.AutoMarshal
	value int32
}] struct{}

var _ __is_ProfileID[ProfileID]

func (x *ProfileID) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ProfileID.WeaverMarshal: nil receiver"))
	}
	enc.Int32(x.value)
}

func (x *ProfileID) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ProfileID.WeaverUnmarshal: nil receiver"))
	}
	x.value = dec.Int32()
}

var _ codegen.AutoMarshal = (*SlotID)(nil)

type __is_SlotID[T ~struct {
	weaver.AutoMarshal
	value int32
}] struct{}

var _ __is_SlotID[SlotID]

func (x *SlotID) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("SlotID.WeaverMarshal: nil receiver"))
	}
	enc.Int32(x.value)
}

func (x *SlotID) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("SlotID.WeaverUnmarshal: nil receiver"))
	}
	x.value = dec.Int32()
}

var _ codegen.AutoMarshal = (*StampCondition)(nil)

type __is_StampCondition[T ~struct {
	weaver.AutoMarshal
	value StampConditionKind
}] struct{}

var _ __is_StampCondition[StampCondition]

func (x *StampCondition) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("StampCondition.WeaverMarshal: nil receiver"))
	}
	enc.String((string)(x.value))
}

func (x *StampCondition) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("StampCondition.WeaverUnmarshal: nil receiver"))
	}
	*(*string)(&x.value) = dec.String()
}

var _ codegen.AutoMarshal = (*TalkID)(nil)

type __is_TalkID[T ~struct {
	weaver.AutoMarshal
	value int32
}] struct{}

var _ __is_TalkID[TalkID]

func (x *TalkID) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TalkID.WeaverMarshal: nil receiver"))
	}
	enc.Int32(x.value)
}

func (x *TalkID) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TalkID.WeaverUnmarshal: nil receiver"))
	}
	x.value = dec.Int32()
}

var _ codegen.AutoMarshal = (*TrackID)(nil)

type __is_TrackID[T ~struct {
	weaver.AutoMarshal
	value int32
}] struct{}

var _ __is_TrackID[TrackID]

func (x *TrackID) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TrackID.WeaverMarshal: nil receiver"))
	}
	enc.Int32(x.value)
}

func (x *TrackID) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TrackID.WeaverUnmarshal: nil receiver"))
	}
	x.value = dec.Int32()
}

var _ codegen.AutoMarshal = (*TrackName)(nil)

type __is_TrackName[T ~struct {
	weaver.AutoMarshal
	value string
}] struct{}

var _ __is_TrackName[TrackName]

func (x *TrackName) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TrackName.WeaverMarshal: nil receiver"))
	}
	enc.String(x.value)
}

func (x *TrackName) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TrackName.WeaverUnmarshal: nil receiver"))
	}
	x.value = dec.String()
}

var _ codegen.AutoMarshal = (*ViewingSeconds)(nil)

type __is_ViewingSeconds[T ~struct {
	weaver.AutoMarshal
	value int32
}] struct{}

var _ __is_ViewingSeconds[ViewingSeconds]

func (x *ViewingSeconds) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ViewingSeconds.WeaverMarshal: nil receiver"))
	}
	enc.Int32(x.value)
}

func (x *ViewingSeconds) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ViewingSeconds.WeaverUnmarshal: nil receiver"))
	}
	x.value = dec.Int32()
}
