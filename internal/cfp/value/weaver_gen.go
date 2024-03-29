// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package value

import (
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

var _ codegen.AutoMarshal = (*SpanSeconds)(nil)

type __is_SpanSeconds[T ~struct {
	weaver.AutoMarshal
	value int
}] struct{}

var _ __is_SpanSeconds[SpanSeconds]

func (x *SpanSeconds) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("SpanSeconds.WeaverMarshal: nil receiver"))
	}
	enc.Int(x.value)
}

func (x *SpanSeconds) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("SpanSeconds.WeaverUnmarshal: nil receiver"))
	}
	x.value = dec.Int()
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

var _ codegen.AutoMarshal = (*VotingTerm)(nil)

type __is_VotingTerm[T ~struct {
	weaver.AutoMarshal
	start time.Time
	end   time.Time
}] struct{}

var _ __is_VotingTerm[VotingTerm]

func (x *VotingTerm) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("VotingTerm.WeaverMarshal: nil receiver"))
	}
	enc.EncodeBinaryMarshaler(&x.start)
	enc.EncodeBinaryMarshaler(&x.end)
}

func (x *VotingTerm) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("VotingTerm.WeaverUnmarshal: nil receiver"))
	}
	dec.DecodeBinaryUnmarshaler(&x.start)
	dec.DecodeBinaryUnmarshaler(&x.end)
}
