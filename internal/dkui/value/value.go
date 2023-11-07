package value

import (
	"fmt"
	"regexp"

	"github.com/ServiceWeaver/weaver"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	TIMEWINDOW_VIEWER_COUNT = 30
	METRICS_UPDATE_INTERVAL = 60
)

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

// ProfileID represents an ID of user profile.
type ProfileID struct {
	weaver.AutoMarshal
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
	weaver.AutoMarshal
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

// SlotID represents an ID of talk slot.
type SlotID struct {
	weaver.AutoMarshal
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
	weaver.AutoMarshal
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
	weaver.AutoMarshal
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
	TALK_SECONDS        = 2400
	STAMP_READY_SECONDS = 1200
	INTERVAL_SECONDS    = 120
	GUARD_SECONDS       = INTERVAL_SECONDS - 10
)

var (
	ViewingSeconds120  ViewingSeconds
	ViewingSeconds2400 ViewingSeconds
)

func init() {
	ViewingSeconds120, _ = NewViewingSeconds(120)
	ViewingSeconds2400, _ = NewViewingSeconds(2400)
}

// ChannelArn represents an ARN of Channel.
type ChannelArn struct {
	weaver.AutoMarshal
	value string
}

func NewChannelArn(v string) (ChannelArn, error) {
	o := ChannelArn{value: v}
	return o, o.Validate()
}

func (v *ChannelArn) Value() arn.ARN {
	a, _ := arn.Parse((v.value))
	return a
}

func (v *ChannelArn) String() string {
	return v.value
}

func (v *ChannelArn) Validate() error {
	if !arn.IsARN(v.value) {
		return fmt.Errorf("is not ARN : %v", v.value)
	}
	if _, err := arn.Parse(v.value); err != nil {
		return fmt.Errorf("value(%v) cannot be parsed: %w", v.value, err)
	}
	return nil
}

// TrackName represents a track name of Dreamkast.
type TrackName struct {
	weaver.AutoMarshal
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
