// go:build !ignoreWeaverGen

package value

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"fmt"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

// Local stub implementations.

// Client stub implementations.

// Server stub implementations.

// AutoMarshal implementations.

var _ codegen.AutoMarshal = &ConfName{}

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

var _ codegen.AutoMarshal = &TalkID{}

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
