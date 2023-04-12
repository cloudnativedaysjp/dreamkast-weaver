package cfp_test

import (
	"context"
	"net"
	"testing"
	"time"

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/sqlhelper"
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := voter.GetCount(ctx, cfp.GetCountRequest{
		ConfName: "cndf2023",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ok bool
	for _, r := range resp {
		if r.TalkID == talkID {
			ok = true
			if r.Count == 0 {
				t.Errorf("count must be more than 0")
			}
		}
	}
	if !ok {
		t.Errorf("talkID not found")
	}

}
