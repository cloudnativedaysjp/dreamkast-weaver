package value

import (
	"errors"
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
		validation.In(cicd2023, cndf2023, cndt2023),
		validation.By(func(value interface{}) error {
			v, ok := value.(ConferenceKind)
			if !ok {
				return errors.New("must be conference kind")
			}
			if _, ok := cfpTerms[v]; !ok {
				return errors.New("must be set to cfp terms")
			}
			return nil
		}),
	)
}

func (v *ConfName) Start() time.Time {
	return cfpTerms[v.value].start
}

func (v *ConfName) End() time.Time {
	return cfpTerms[v.value].end
}

const (
	SPAN_SECONDS = 3600
)

type ConferenceKind string
type cfpTerm struct {
	start time.Time
	end   time.Time
}

var (
	cicd2023 ConferenceKind = "cicd2023"
	cndf2023 ConferenceKind = "cndf2023"
	cndt2023 ConferenceKind = "cndt2023"

	CICD2023 ConfName
	CNDF2023 ConfName
	CNDT2023 ConfName

	jst *time.Location

	cfpTerms map[ConferenceKind]cfpTerm
)

func init() {
	var err, e error
	jst, e = time.LoadLocation("Asia/Tokyo")
	err = errors.Join(err, e)

	cfpTerms = make(map[ConferenceKind]cfpTerm)
	cfpTerms[cicd2023] = cfpTerm{
		start: time.Date(2023, 1, 1, 0, 0, 0, 0, jst),
		end:   time.Date(2023, 1, 25, 18, 0, 0, 0, jst),
	}
	cfpTerms[cndf2023] = cfpTerm{
		start: time.Date(2023, 5, 2, 0, 0, 0, 0, jst),
		end:   time.Date(2023, 6, 25, 18, 0, 0, 0, jst), // TODO adjust
	}
	cfpTerms[cndt2023] = cfpTerm{
		start: time.Date(2023, 9, 1, 0, 0, 0, 0, jst),    // TODO adjust
		end:   time.Date(2023, 11, 25, 18, 0, 0, 0, jst), // TODO adjust
	}

	CICD2023, e = NewConfName(cicd2023)
	err = errors.Join(err, e)
	CNDF2023, e = NewConfName(cndf2023)
	err = errors.Join(err, e)
	CNDT2023, e = NewConfName(cndt2023)
	err = errors.Join(err, e)
	if err != nil {
		panic(err)
	}
}
