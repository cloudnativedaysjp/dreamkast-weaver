package domain

import (
	"dreamkast-weaver/internal/cfp/value"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type CfpDomain struct{}

type CfpVote struct {
	weaver.AutoMarshal
	TalkID    value.TalkID
	GlobalIP  value.GlobalIP
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
			ip:        v.GlobalIP.Value(),
			timeFrame: v.CreatedAt.Unix() / value.SPAN_SECONDS,
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
