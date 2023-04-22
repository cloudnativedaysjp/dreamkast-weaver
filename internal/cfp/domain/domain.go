package domain

import (
	"dreamkast-weaver/internal/cfp/value"

	"github.com/ServiceWeaver/weaver"
)

type CfpDomain struct{}

type CfpVote struct {
	weaver.AutoMarshal
	TalkID   value.TalkID
	GlobalIP value.GlobalIP
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

func (cv *CfpDomain) TallyCfpVotes(cvs *CfpVotes) ([]*VoteCount, error) {
	counts := map[int32]int{}
	for _, cv := range cvs.Items {
		counts[cv.TalkID.Value()]++
	}

	var resp []*VoteCount
	for talkID, count := range counts {
		talkID, err := value.NewTalkID(talkID)
		if err != nil {
			return nil, err
		}
		resp = append(resp, &VoteCount{
			TalkID: talkID,
			Count:  count,
		})
	}
	return resp, nil
}
