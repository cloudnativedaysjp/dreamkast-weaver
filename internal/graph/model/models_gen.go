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
	ProfileID int      `json:"profileID"`
	TrackID   int      `json:"trackID"`
	TalkID    int      `json:"talkID"`
	SlotID    int      `json:"slotID"`
}

type StampChallenge struct {
	SlotID    int                `json:"slotID"`
	Condition ChallengeCondition `json:"condition"`
	UpdatedAt int                `json:"updatedAt"`
}

type StampOnSiteInput struct {
	ConfName  ConfName `json:"confName"`
	ProfileID int      `json:"profileID"`
	TrackID   int      `json:"trackID"`
	TalkID    int      `json:"talkID"`
	SlotID    int      `json:"slotID"`
}

type StampOnlineInput struct {
	ConfName  ConfName `json:"confName"`
	ProfileID int      `json:"profileID"`
	SlotID    int      `json:"slotID"`
}

type ViewTrackInput struct {
	ProfileID int    `json:"profileID"`
	TrackName string `json:"trackName"`
}

type ViewerCount struct {
	TrackName string `json:"trackName"`
	Count     int    `json:"count"`
}

type ViewingSlot struct {
	SlotID      int `json:"slotId"`
	ViewingTime int `json:"viewingTime"`
}

type VoteCount struct {
	TalkID int `json:"talkId"`
	Count  int `json:"count"`
}

type VoteInput struct {
	ConfName ConfName `json:"confName"`
	TalkID   int      `json:"talkId"`
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

func (e *ChallengeCondition) UnmarshalGQL(v interface{}) error {
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
)

var AllConfName = []ConfName{
	ConfNameCicd2023,
	ConfNameCndf2023,
	ConfNameCndt2023,
}

func (e ConfName) IsValid() bool {
	switch e {
	case ConfNameCicd2023, ConfNameCndf2023, ConfNameCndt2023:
		return true
	}
	return false
}

func (e ConfName) String() string {
	return string(e)
}

func (e *ConfName) UnmarshalGQL(v interface{}) error {
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
