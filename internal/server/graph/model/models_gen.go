// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type CreateViewEventInput struct {
	ConfName  ConfName `json:"confName"`
	ProfileID int32    `json:"profileID"`
	TrackID   int32    `json:"trackID"`
	TalkID    int32    `json:"talkID"`
	SlotID    int32    `json:"slotID"`
}

type Mutation struct {
}

type Query struct {
}

type StampChallenge struct {
	SlotID    int32              `json:"slotID"`
	Condition ChallengeCondition `json:"condition"`
	UpdatedAt int32              `json:"updatedAt"`
}

type StampOnSiteInput struct {
	ConfName  ConfName `json:"confName"`
	ProfileID int32    `json:"profileID"`
	TrackID   int32    `json:"trackID"`
	TalkID    int32    `json:"talkID"`
	SlotID    int32    `json:"slotID"`
}

type StampOnlineInput struct {
	ConfName  ConfName `json:"confName"`
	ProfileID int32    `json:"profileID"`
	SlotID    int32    `json:"slotID"`
}

type ViewTrackInput struct {
	ProfileID int32  `json:"profileID"`
	TrackName string `json:"trackName"`
	TalkID    int32  `json:"talkID"`
}

type ViewerCount struct {
	TrackName string `json:"trackName"`
	Count     int32  `json:"count"`
}

type ViewingSlot struct {
	SlotID      int32 `json:"slotId"`
	ViewingTime int32 `json:"viewingTime"`
}

type VoteCount struct {
	TalkID int32 `json:"talkId"`
	Count  int32 `json:"count"`
}

type VoteInput struct {
	ConfName ConfName `json:"confName"`
	TalkID   int32    `json:"talkId"`
}

type VotingTerm struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

type ChallengeCondition string

const (
	ChallengeConditionReady   ChallengeCondition = "READY"
	ChallengeConditionStamped ChallengeCondition = "STAMPED"
	ChallengeConditionSkipped ChallengeCondition = "SKIPPED"
)

var AllChallengeCondition = []ChallengeCondition{
	ChallengeConditionReady,
	ChallengeConditionStamped,
	ChallengeConditionSkipped,
}

func (e ChallengeCondition) IsValid() bool {
	switch e {
	case ChallengeConditionReady, ChallengeConditionStamped, ChallengeConditionSkipped:
		return true
	}
	return false
}

func (e ChallengeCondition) String() string {
	return string(e)
}

func (e *ChallengeCondition) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChallengeCondition(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChallengeCondition", str)
	}
	return nil
}

func (e ChallengeCondition) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ConfName string

const (
	ConfNameCicd2023 ConfName = "cicd2023"
	ConfNameCndf2023 ConfName = "cndf2023"
	ConfNameCndt2023 ConfName = "cndt2023"
	ConfNameCnds2024 ConfName = "cnds2024"
	ConfNameCndw2024 ConfName = "cndw2024"
	ConfNameCnds2025 ConfName = "cnds2025"
)

var AllConfName = []ConfName{
	ConfNameCicd2023,
	ConfNameCndf2023,
	ConfNameCndt2023,
	ConfNameCnds2024,
	ConfNameCndw2024,
	ConfNameCnds2025,
}

func (e ConfName) IsValid() bool {
	switch e {
	case ConfNameCicd2023, ConfNameCndf2023, ConfNameCndt2023, ConfNameCnds2024, ConfNameCndw2024, ConfNameCnds2025:
		return true
	}
	return false
}

func (e ConfName) String() string {
	return string(e)
}

func (e *ConfName) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ConfName(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ConfName", str)
	}
	return nil
}

func (e ConfName) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
