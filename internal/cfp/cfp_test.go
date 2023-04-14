package cfp_test

import (
	"context"
	"net"
	"testing"
	"time"

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/sqlhelper"

	"github.com/stretchr/testify/assert"
)

func TestCfpVoteImpl(t *testing.T) {
	sh := sqlhelper.NewTestSqlHelper()
	voter := cfp.NewVoter(sh)
	ctx := context.Background()

	talkID := int32(time.Now().Unix())

	err := voter.Vote(ctx, cfp.VoteRequest{
		ConfName: "cndf2023",
		TalkID:   talkID,
		GlobalIP: net.ParseIP("127.0.0.1"),
	})
	assert.Nil(t, err)

	resp, err := voter.GetCount(ctx, cfp.GetCountRequest{
		ConfName: "cndf2023",
	})
	assert.Nil(t, err)

	var ok bool
	for _, r := range resp {
		if r.TalkID == talkID {
			ok = true
			assert.Greater(t, r.Count, 0)
		}
	}
	assert.True(t, ok, "talkID not found")

}
