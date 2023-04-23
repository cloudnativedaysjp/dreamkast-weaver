package domain_test

import (
	"dreamkast-weaver/internal/cfp/domain"
	"dreamkast-weaver/internal/cfp/value"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	svc = domain.CfpDomain{}
)

func TestCfpDomain_TallyCfpVotes(t *testing.T) {

	tn := time.Unix(time.Now().Unix()/int64(value.SPAN_SECONDS)*int64(value.SPAN_SECONDS), 0)
	id := newTalkID(1)
	//ip := NewGlobalIP("192.0.2.1")
	ip := NewGlobalIP("192.168.100.1")

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
						GlobalIP:  ip,
						CreatedAt: tn,
					},
					{
						TalkID:    id,
						GlobalIP:  ip,
						CreatedAt: tn.Add((value.SPAN_SECONDS - 1) * time.Second),
					},
					{
						TalkID:    id,
						GlobalIP:  ip,
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
						GlobalIP:  ip,
						CreatedAt: tn,
					},
					{
						TalkID:    newTalkID(2),
						GlobalIP:  ip,
						CreatedAt: tn.Add((value.SPAN_SECONDS - 3) * time.Second),
					},
					{
						TalkID:    id,
						GlobalIP:  NewGlobalIP("192.0.2.2"),
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

			got := svc.TallyCfpVotes(cvs)

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

func NewGlobalIP(v string) value.GlobalIP {
	ip, err := value.NewGlobalIP(v)
	mustNil(err)
	return ip
}

func newTalkID(v int32) value.TalkID {
	id, err := value.NewTalkID(v)
	mustNil(err)
	return id
}
