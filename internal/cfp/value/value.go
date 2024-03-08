package value

import (
	"fmt"
	"time"

	"github.com/ServiceWeaver/weaver"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

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
		validation.In(
			cicd2023, cndf2023, cndt2023,
			cnds2024,
		),
	)
}

type ConferenceKind string

var (
	cicd2023 ConferenceKind = "cicd2023"
	cndf2023 ConferenceKind = "cndf2023"
	cndt2023 ConferenceKind = "cndt2023"
	cnds2024 ConferenceKind = "cnds2024"

	CICD2023 ConfName
	CNDF2023 ConfName
	CNDT2023 ConfName
	CNDS2024 ConfName
)

func init() {
	CICD2023, _ = NewConfName(cicd2023)
	CNDF2023, _ = NewConfName(cndf2023)
	CNDT2023, _ = NewConfName(cndt2023)
	CNDS2024, _ = NewConfName(cnds2024)
}

var (
	jst               *time.Location
	defaultVotingTerm VotingTerm
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	defaultVotingTerm = VotingTerm{
		start: time.Date(2000, 1, 1, 0, 0, 0, 0, jst),
		end:   time.Date(2100, 12, 31, 23, 59, 59, 0, jst),
	}
}

type VotingTerm struct {
	weaver.AutoMarshal
	start time.Time
	end   time.Time
}

func NewVotingTerm(s, e *time.Time) (VotingTerm, error) {
	o := VotingTerm{
		start: defaultVotingTerm.start,
		end:   defaultVotingTerm.end,
	}

	if s != nil {
		o.start = s.In(jst)
	}

	if e != nil {
		o.end = e.In(jst)
	}

	return o, o.Validate()
}

func (v *VotingTerm) Value() (time.Time, time.Time) {
	return v.start, v.end
}

func (v *VotingTerm) String() string {
	return fmt.Sprintf("%v-%v", v.start, v.end)
}

func (v *VotingTerm) Validate() error {
	if v.start.After(v.end) {
		return fmt.Errorf("start(%v) should be before end(%v)", v.start, v.end)
	}

	return nil
}

type SpanSeconds struct {
	weaver.AutoMarshal
	value int
}

const (
	SPAN_SECONDS = 3600
)

func NewSpanSeconds(ss *int) (SpanSeconds, error) {
	o := SpanSeconds{
		value: SPAN_SECONDS,
	}

	if ss != nil {
		o.value = *ss
	}

	return o, o.Validate()
}

func (v *SpanSeconds) Value() int {
	return v.value
}

func (v *SpanSeconds) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *SpanSeconds) Validate() error {
	if v.value < 1 {
		return fmt.Errorf("value(%v) must be greater than or equal to 1", v.value)
	}

	// KIME no atai
	if v.value > 3600*24 {
		return fmt.Errorf("value(%v) must be less than or equal to 1 day", v.value)
	}

	return nil
}
