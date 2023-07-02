package aws

import (
	"context"
	"dreamkast-weaver/internal/dkui/value"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
)

type AWSIVSClient interface {
	GetStream(ctx context.Context, ca value.ChannelArn) (*types.Stream, error)
}

type AWSIVSClientImpl struct {
	client *ivs.Client
}

var _ AWSIVSClient = (*AWSIVSClientImpl)(nil)

func NewAWSIVSClientImpl() (AWSIVSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &AWSIVSClientImpl{
		client: ivs.NewFromConfig(cfg),
	}, nil
}

func (c *AWSIVSClientImpl) GetStream(ctx context.Context, ca value.ChannelArn) (*types.Stream, error) {
	v := ca.String()
	o, err := c.client.GetStream(ctx, &ivs.GetStreamInput{
		ChannelArn: &v,
	})
	if err != nil {
		return nil, err
	}

	return o.Stream, nil
}
