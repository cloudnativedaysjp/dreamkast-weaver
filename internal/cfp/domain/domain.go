package domain

import (
	"dreamkast-weaver/internal/cfp/value"
	"dreamkast-weaver/internal/stacktrace"
	"errors"
	"net"
	"time"

	"github.com/ServiceWeaver/weaver"
)

const (
	SPAN_SECONDS = 3600
)

var (
	jst         *time.Location
	votingTerms map[value.ConfName]VotingTerm
)

var (
	ErrMissingVotingTerm = errors.New("missing voting term")
)

func init() {
	var err error
	jst, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	votingTerms = make(map[value.ConfName]VotingTerm)
	votingTerms[value.CICD2023] = VotingTerm{
		Start: time.Date(2023, 1, 1, 0, 0, 0, 0, jst),
		End:   time.Date(2023, 1, 25, 18, 0, 0, 0, jst),
	}
	votingTerms[value.CNDF2023] = VotingTerm{
		Start: time.Date(2023, 5, 2, 0, 0, 0, 0, jst),
		End:   time.Date(2023, 6, 19, 23, 59, 0, 0, jst),
	}
	votingTerms[value.CNDT2023] = VotingTerm{
		Start: time.Date(2023, 9, 1, 0, 0, 0, 0, jst),    // TODO adjust
		End:   time.Date(2023, 11, 25, 18, 0, 0, 0, jst), // TODO adjust
	}
}

type CfpDomain struct{}

type CfpVote struct {
	weaver.AutoMarshal
	TalkID    value.TalkID
	ClientIp  net.IP
	CreatedAt time.Time
}

type CfpVotes struct {
	weaver.AutoMarshal
	Items []CfpVote
}

type VoteCount struct {
	weaver.AutoMarshal
	TalkID value.TalkID
	Count  int
}

type VotingTerm struct {
	Start time.Time
	End   time.Time
}

func (cd *CfpDomain) TallyCfpVotes(cfpVotes *CfpVotes) []*VoteCount {
	type key struct {
		talkId    int32
		ip        string
		timeFrame int64
	}
	voted := map[key]bool{}
	counts := map[value.TalkID]int{}
	for _, v := range cfpVotes.Items {
		k := key{
			talkId:    v.TalkID.Value(),
			ip:        v.ClientIp.String(),
			timeFrame: v.CreatedAt.Unix() / SPAN_SECONDS,
		}
		if _, isThere := voted[k]; isThere {
			continue
		}
		voted[k] = true
		counts[v.TalkID]++
	}

	var resp []*VoteCount
	for talkID, count := range counts {
		resp = append(resp, &VoteCount{
			TalkID: talkID,
			Count:  count,
		})
	}
	return resp
}

func (cd *CfpDomain) GetVotingTerm(cn value.ConfName) (*VotingTerm, error) {
	if vt, ok := votingTerms[cn]; ok {
		return &vt, nil
	}
	return nil, stacktrace.With(ErrMissingVotingTerm)
}
