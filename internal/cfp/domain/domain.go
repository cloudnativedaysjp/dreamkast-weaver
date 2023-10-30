package domain

import (
	"net"
	"time"

	"github.com/ServiceWeaver/weaver"

	"dreamkast-weaver/internal/cfp/value"
)

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

type CfpDomain struct{}

func (cd *CfpDomain) TallyCfpVotes(cfpVotes *CfpVotes, spanSeconds value.SpanSeconds) []*VoteCount {
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
			timeFrame: v.CreatedAt.Unix() / int64(spanSeconds.Value()),
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
