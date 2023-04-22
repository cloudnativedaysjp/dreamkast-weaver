package value

import (
	"fmt"

	"github.com/ServiceWeaver/weaver"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// GlobalIP represents global ip address.
type GlobalIP struct {
	weaver.AutoMarshal
	value string
}

func NewGlobalIP(v string) (GlobalIP, error) {
	o := GlobalIP{value: v}
	return o, o.Validate()
}

func (v *GlobalIP) Validate() error {
	return validation.Validate(v.value, is.IP)
}

func (v *GlobalIP) Value() string {
	return v.value
}

// TalkID represents an ID of talk.
type TalkID struct {
	weaver.AutoMarshal
	value int32
}

func NewTalkID(v int32) (TalkID, error) {
	o := TalkID{value: v}
	return o, o.Validate()
}

func (v *TalkID) Value() int32 {
	return v.value
}

func (v *TalkID) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *TalkID) Validate() error {
	return nil
}

// ConfName represents a conference name.
type ConfName struct {
	weaver.AutoMarshal
	value ConferenceKind
}

func NewConfName(v ConferenceKind) (ConfName, error) {
	o := ConfName{value: v}
	return o, o.Validate()
}

func (v *ConfName) Value() ConferenceKind {
	return v.value
}

func (v *ConfName) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *ConfName) Validate() error {
	return validation.Validate(v.value,
		validation.In(cicd2023, cndf2023, cndt2023),
	)
}

type ConferenceKind string

var (
	cicd2023 ConferenceKind = "cicd2023"
	cndf2023 ConferenceKind = "cndf2023"
	cndt2023 ConferenceKind = "cndt2023"

	CICD2023 ConfName
	CNDF2023 ConfName
	CNDT2023 ConfName
)

func init() {
	CICD2023, _ = NewConfName(cicd2023)
	CNDF2023, _ = NewConfName(cndf2023)
	CNDT2023, _ = NewConfName(cndt2023)
}
