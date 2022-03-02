package cloudpubsub

import (
	"context"
	"encoding/json"
	"time"

	pubsub "cloud.google.com/go/pubsub"
	pubsubapi "cloud.google.com/go/pubsub/apiv1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/rabee-inc/go-pkg/log"
)

// Client ... PubSubのクライアント
type Client struct {
	cPubSub     *pubsub.Client
	cSubscriber *pubsubapi.SubscriberClient
	projectID   string
}

// NewClient ... PubSubのクライアントを取得する
func NewClient(projectID string) *Client {
	ctx := context.Background()
	gOpt := option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                1 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}))
	cPubSub, err := pubsub.NewClient(ctx, projectID, gOpt)
	if err != nil {
		panic(err)
	}
	cSubscriber, err := pubsubapi.NewSubscriberClient(ctx, gOpt)
	if err != nil {
		panic(err)
	}
	return &Client{
		cPubSub,
		cSubscriber,
		projectID,
	}
}

// Publish ... メッセージを送信する
func (c *Client) Publish(
	ctx context.Context,
	topicID string,
	msg interface{},
) error {
	bMsg, err := json.Marshal(msg)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	if _, err := c.cPubSub.Topic(topicID).Publish(ctx, &pubsub.Message{
		Data: bMsg,
	}).Get(ctx); err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}
