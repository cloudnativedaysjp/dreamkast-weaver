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
	sqlh := sqlhelper.NewTestSqlHelper()
	ctx := sqlhelper.WithSqlHelper(context.Background(), sqlh)

	talkID := int32(time.Now().Unix())

	err := cfp.NewCfpVote().Vote(ctx, cfp.VoteRequest{
		ConfName: "cndf2023",
		TalkID:   talkID,
		GlobalIP: net.ParseIP("127.0.0.1"),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := cfp.NewCfpVote().GetCount(ctx, cfp.GetCountRequest{
		ConfName: "cndf2023",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ok bool
	for _, r := range resp {
		if r.TalkId == talkID {
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
