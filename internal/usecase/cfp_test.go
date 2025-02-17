package usecase

import (
	"context"
	"net"
	"testing"
	"time"

	"dreamkast-weaver/internal/domain/value"
	"dreamkast-weaver/internal/sqlhelper"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestCfpVoteImpl(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	cn := value.CICD2023
	talkID, _ := value.NewTalkID(3)

	sq, _ := sqlhelper.NewSqlHelper(&sqlhelper.SqlOption{
		User:     "user",
		Password: "password",
		Endpoint: "127.0.0.1",
		Port:     "13306",
		DbName:   "test_cfp",
	})

	svc := NewCFPService(sq)

	err := svc.Vote(ctx, VoteRequest{
		ConfName: cn,
		TalkID:   talkID,
		ClientIp: net.ParseIP("127.0.0.1"),
	})
	assert.NoError(t, err)

	vts := time.Now().AddDate(0, 0, -1)
	vte := time.Now().AddDate(0, 0, 1)
	vt, _ := value.NewVotingTerm(&vts, &vte)

	resp, err := svc.VoteCounts(ctx, VoteCountsRequest{
		ConfName:    cn,
		VotingTerm:  vt,
		SpanSeconds: newSpanSeconds(value.SPAN_SECONDS),
	})
	assert.NoError(t, err)

	var ok bool
	for _, r := range resp {
		if r.TalkID == talkID {
			ok = true
			assert.Greater(t, r.Count, 0)
		}
	}
	assert.True(t, ok, "talkID not found")
}

func newSpanSeconds(v int) value.SpanSeconds {
	ss, err := value.NewSpanSeconds(&v)
	mustNil(err)
	return ss
}
