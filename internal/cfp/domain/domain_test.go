package domain_test

import (
	"dreamkast-weaver/internal/cfp/domain"
	"dreamkast-weaver/internal/cfp/value"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	svc = domain.CfpDomain{}
)

func TestCfpDomain_TallyCfpVotes(t *testing.T) {
	ss := newSpanSeconds(value.SPAN_SECONDS)
	tn := time.Unix(time.Now().Unix()/value.SPAN_SECONDS*value.SPAN_SECONDS, 0)
	id := newTalkID(1)
	ip := net.ParseIP("192.0.2.1")

	tests := []struct {
		name   string
		given  func() (cvs *domain.CfpVotes)
		counts map[value.TalkID]int
	}{
		{
			name: "votes within a span are summarized",
			given: func() (cvs *domain.CfpVotes) {
				cvs = &domain.CfpVotes{Items: []domain.CfpVote{
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn,
					},
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn.Add((value.SPAN_SECONDS - 1) * time.Second),
					},
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn.Add(value.SPAN_SECONDS * time.Second),
					},
				}}
				return cvs
			},
			counts: map[value.TalkID]int{
				id: 2,
			},
		},
		{
			name: "different IPs or IDs are different votes",
			given: func() (cvs *domain.CfpVotes) {
				cvs = &domain.CfpVotes{Items: []domain.CfpVote{
					{
						TalkID:    id,
						ClientIp:  ip,
						CreatedAt: tn,
					},
					{
						TalkID:    newTalkID(2),
						ClientIp:  ip,
						CreatedAt: tn.Add((value.SPAN_SECONDS - 3) * time.Second),
					},
					{
						TalkID:    id,
						ClientIp:  net.ParseIP("192.0.2.2"),
						CreatedAt: tn.Add((value.SPAN_SECONDS - 1) * time.Second),
					},
				}}
				return cvs
			},
			counts: map[value.TalkID]int{
				id:           2,
				newTalkID(2): 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run("ok:"+tt.name, func(t *testing.T) {
			cvs := tt.given()

			got := svc.TallyCfpVotes(cvs, ss)

			for _, v := range got {
				assert.Equal(t, tt.counts[v.TalkID], v.Count)
			}
		})
	}
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}

func newTalkID(v int32) value.TalkID {
	id, err := value.NewTalkID(v)
	mustNil(err)
	return id
}

func newSpanSeconds(v int) value.SpanSeconds {
	ss, err := value.NewSpanSeconds(&v)
	mustNil(err)
	return ss
}
