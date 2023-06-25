package aws

import (
	"context"
	"dreamkast-weaver/internal/dkui/value"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
)

type AWSClient interface {
	IVSGetStream(ctx context.Context, ca value.ChannelArn) (*types.Stream, error)
}

type AWSClientImpl struct {
	ivsClient *ivs.Client
}

var _ AWSClient = (*AWSClientImpl)(nil)

func NewAWSClientImpl() (AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &AWSClientImpl{
		ivsClient: ivs.NewFromConfig(cfg),
	}, nil
}

func (c *AWSClientImpl) IVSGetStream(ctx context.Context, ca value.ChannelArn) (*types.Stream, error) {
	v := ca.String()
	o, err := c.ivsClient.GetStream(ctx, &ivs.GetStreamInput{
		ChannelArn: &v,
	})
	if err != nil {
		return nil, err
	}

	return o.Stream, nil
}
