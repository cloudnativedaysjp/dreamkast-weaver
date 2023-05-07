package domain

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"fmt"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

// Local stub implementations.

// Client stub implementations.

// Server stub implementations.

// AutoMarshal implementations.

var _ codegen.AutoMarshal = &CfpVote{}

func (x *CfpVote) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("CfpVote.WeaverMarshal: nil receiver"))
	}
	(x.TalkID).WeaverMarshal(enc)
	serviceweaver_enc_slice_byte_87461245(enc, ([]byte)(x.ClientIp))
	enc.EncodeBinaryMarshaler(&x.CreatedAt)
}

func (x *CfpVote) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("CfpVote.WeaverUnmarshal: nil receiver"))
	}
	(&x.TalkID).WeaverUnmarshal(dec)
	*(*[]byte)(&x.ClientIp) = serviceweaver_dec_slice_byte_87461245(dec)
	dec.DecodeBinaryUnmarshaler(&x.CreatedAt)
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

var _ codegen.AutoMarshal = &CfpVotes{}

func (x *CfpVotes) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("CfpVotes.WeaverMarshal: nil receiver"))
	}
	serviceweaver_enc_slice_CfpVote_9206e939(enc, x.Items)
}

func (x *CfpVotes) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("CfpVotes.WeaverUnmarshal: nil receiver"))
	}
	x.Items = serviceweaver_dec_slice_CfpVote_9206e939(dec)
}

func serviceweaver_enc_slice_CfpVote_9206e939(enc *codegen.Encoder, arg []CfpVote) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_CfpVote_9206e939(dec *codegen.Decoder) []CfpVote {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]CfpVote, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}

var _ codegen.AutoMarshal = &VoteCount{}

func (x *VoteCount) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("VoteCount.WeaverMarshal: nil receiver"))
	}
	(x.TalkID).WeaverMarshal(enc)
	enc.Int(x.Count)
}

func (x *VoteCount) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("VoteCount.WeaverUnmarshal: nil receiver"))
	}
	(&x.TalkID).WeaverUnmarshal(dec)
	x.Count = dec.Int()
}
