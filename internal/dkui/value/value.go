package value

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ConfName represents a conference name.
type ConfName struct {
	valueObject[ConferenceKind]
}

func NewConfName(v ConferenceKind) (ConfName, error) {
	o := ConfName{wrap(v)}
	return o, o.Validate()
}

func (v *ConfName) Validate() error {
	return validation.Validate(v.Value(),
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
	valueObject[int32]
}

func NewProfileID(v int32) (ProfileID, error) {
	return ProfileID{wrap(v)}, nil
}

// TrackID represents an ID of talk track.
type TrackID struct {
	valueObject[int32]
}

func NewTrackID(v int32) (TrackID, error) {
	return TrackID{wrap(v)}, nil
}

// TrackID represents an ID of talk.
type TalkID struct {
	valueObject[int32]
}

func NewTalkID(v int32) (TalkID, error) {
	return TalkID{wrap(v)}, nil
}

// SlotID represents an ID of talk slot.
type SlotID struct {
	valueObject[int32]
}

func NewSlotID(v int32) (SlotID, error) {
	return SlotID{wrap(v)}, nil
}

// StampCondition represents a condition of the stamp of talk slot.
type StampCondition struct {
	valueObject[StampConditionKind]
}

func NewStampCondtion(v StampConditionKind) (StampCondition, error) {
	o := StampCondition{wrap(v)}
	return o, o.Validate()
}

func (v *StampCondition) Validate() error {
	return validation.Validate(v.Value(),
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
	StampReady, _ = NewStampCondtion(stampReady)
	StampStamped, _ = NewStampCondtion(stampStamped)
	StampSkipped, _ = NewStampCondtion(stampSkipped)
}

// ViewingSeconds represents a talk viewing seconds.
type ViewingSeconds struct {
	valueObject[int32]
}

func NewViewingPeriod(v int32) (ViewingSeconds, error) {
	o := ViewingSeconds{wrap(v)}
	return o, o.Validate()
}

func (v *ViewingSeconds) Validate() error {
	return validation.Validate(v.Value(),
		validation.In(INTERVAL_SECONDS, TALK_SECONDS),
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
	ViewingSeconds120, _ = NewViewingPeriod(120)
	ViewingSeconds2400, _ = NewViewingPeriod(2400)
}
