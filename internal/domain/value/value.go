package value

import (
	"fmt"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	TIMEWINDOW_VIEWER_COUNT = 60
	METRICS_UPDATE_INTERVAL = 30
)

// ConfName represents a conference name.
type ConfName struct {
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
			cnds2024, cndw2024,
			cnds2025,
		),
	)
}

type ConferenceKind string

var (
	cicd2023 ConferenceKind = "cicd2023"
	cndf2023 ConferenceKind = "cndf2023"
	cndt2023 ConferenceKind = "cndt2023"
	cnds2024 ConferenceKind = "cnds2024"
	cndw2024 ConferenceKind = "cndw2024"
	cnds2025 ConferenceKind = "cnds2025"

	CICD2023 ConfName
	CNDF2023 ConfName
	CNDT2023 ConfName
	CNDS2024 ConfName
	CNDW2024 ConfName
	CNDS2025 ConfName
)

func init() {
	CICD2023, _ = NewConfName(cicd2023)
	CNDF2023, _ = NewConfName(cndf2023)
	CNDT2023, _ = NewConfName(cndt2023)
	CNDS2024, _ = NewConfName(cnds2024)
	CNDW2024, _ = NewConfName(cndw2024)
	CNDS2025, _ = NewConfName(cnds2025)
}

// ProfileID represents an ID of user profile.
type ProfileID struct {
	value int32
}

func NewProfileID(v int32) (ProfileID, error) {
	o := ProfileID{value: v}
	return o, o.Validate()
}

func (v *ProfileID) Value() int32 {
	return v.value
}

func (v *ProfileID) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *ProfileID) Validate() error {
	return nil
}

// TrackID represents an ID of talk track.
type TrackID struct {
	value int32
}

func NewTrackID(v int32) (TrackID, error) {
	o := TrackID{value: v}
	return o, o.Validate()
}

func (v *TrackID) Value() int32 {
	return v.value
}

func (v *TrackID) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *TrackID) Validate() error {
	return nil
}

// TalkID represents an ID of talk.
type TalkID struct {
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

// SlotID represents an ID of talk slot.
type SlotID struct {
	value int32
}

func NewSlotID(v int32) (SlotID, error) {
	o := SlotID{value: v}
	return o, o.Validate()
}

func (v *SlotID) Value() int32 {
	return v.value
}

func (v *SlotID) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *SlotID) Validate() error {
	return nil
}

// StampCondition represents a condition of the stamp of talk slot.
type StampCondition struct {
	value StampConditionKind
}

func NewStampCondition(v StampConditionKind) (StampCondition, error) {
	o := StampCondition{value: v}
	return o, o.Validate()
}

func (v *StampCondition) Value() StampConditionKind {
	return v.value
}

func (v *StampCondition) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *StampCondition) Validate() error {
	return validation.Validate(v.value,
		validation.In(stampReady, stampStamped, stampSkipped),
	)
}

type StampConditionKind string

var (
	stampReady   StampConditionKind = "ready"
	stampStamped StampConditionKind = "stamped"
	stampSkipped StampConditionKind = "skipped"

	StampReady   StampCondition
	StampStamped StampCondition
	StampSkipped StampCondition
)

func init() {
	StampReady, _ = NewStampCondition(stampReady)
	StampStamped, _ = NewStampCondition(stampStamped)
	StampSkipped, _ = NewStampCondition(stampSkipped)
}

// ViewingSeconds represents a talk viewing seconds.
type ViewingSeconds struct {
	value int32
}

func NewViewingSeconds(v int32) (ViewingSeconds, error) {
	o := ViewingSeconds{value: v}
	return o, o.Validate()
}

func (v *ViewingSeconds) Value() int32 {
	return v.value
}

func (v *ViewingSeconds) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *ViewingSeconds) Validate() error {
	return validation.Validate(v.value,
		validation.In(int32(INTERVAL_SECONDS), int32(TALK_SECONDS)),
	)
}

const (
	TALK_SECONDS              = 2400
	STAMP_READY_SECONDS int32 = 1200
	INTERVAL_SECONDS          = 120
	GUARD_SECONDS             = INTERVAL_SECONDS - 10
)

var (
	ViewingSeconds120  ViewingSeconds
	ViewingSeconds2400 ViewingSeconds
)

func init() {
	ViewingSeconds120, _ = NewViewingSeconds(120)
	ViewingSeconds2400, _ = NewViewingSeconds(2400)
}

// TrackName represents a track name of Dreamkast.
type TrackName struct {
	value string
}

func NewTrackName(v string) (TrackName, error) {
	o := TrackName{value: v}
	return o, o.Validate()
}

func (v *TrackName) Value() string {
	return v.value
}

func (v *TrackName) String() string {
	return v.value
}

func (v *TrackName) Validate() error {
	return validation.Validate(v.value,
		// dreamkast has a maximum of 6 tracks
		validation.Match(regexp.MustCompile("^[A-F]$")),
	)
}

func TrackNames() []TrackName {
	var names []TrackName
	for i := 0; i < 6; i++ {
		names = append(names, TrackName{value: string(rune('A' + i))})
	}
	return names
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
	value int
}

const (
	SPAN_SECONDS = 3600
)

func NewSpanSeconds(ss *int32) (SpanSeconds, error) {
	o := SpanSeconds{
		value: SPAN_SECONDS,
	}

	if ss != nil {
		o.value = int(*ss)
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
