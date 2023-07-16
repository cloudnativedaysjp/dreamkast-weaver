package cfp_test

import (
	"context"
	"net"
	"net/url"
	"testing"

	"dreamkast-weaver/internal/cfp"
	"dreamkast-weaver/internal/cfp/value"

	"github.com/ServiceWeaver/weaver/weavertest"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	weaverConfig = `
	["dreamkast-weaver/internal/cfp/Service"]
	db_user = "user"
	db_password = "password"
	db_endpoint = "127.0.0.1"
	db_port = "13306"
	db_name = "test_cfp"
	`
	dbUrl = "mysql://user:password@127.0.0.1:13306/test_cfp"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	u, _ := url.Parse(dbUrl)
	db := dbmate.New(u)
	db.AutoDumpSchema = false

	mustNil(db.Drop())
	mustNil(db.CreateAndMigrate())
}

func teardown() {}

func TestCfpVoteImpl(t *testing.T) {
	t.Skip()

	runner := weavertest.Local
	runner.Config = weaverConfig
	runner.Test(t, func(t *testing.T, svc cfp.Service) {
		ctx := context.Background()
		cn := value.CICD2023
		talkID, _ := value.NewTalkID(3)

		err := svc.Vote(ctx, cfp.VoteRequest{
			ConfName: cn,
			TalkID:   talkID,
			ClientIp: net.ParseIP("127.0.0.1"),
		})
		assert.NoError(t, err)

		resp, err := svc.VoteCounts(ctx, cn)
		assert.NoError(t, err)

		var ok bool
		for _, r := range resp {
			if r.TalkID == talkID {
				ok = true
				assert.Greater(t, r.Count, 0)
			}
		}
		assert.True(t, ok, "talkID not found")
	})
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}
