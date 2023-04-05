package cloudpubsub

import (
	"context"
	"encoding/json"

	pubsub "cloud.google.com/go/pubsub"
	pubsubapi "cloud.google.com/go/pubsub/apiv1"
	"github.com/rabee-inc/go-pkg/log"
)

type Client struct {
	cPubSub     *pubsub.Client
	cSubscriber *pubsubapi.SubscriberClient
	projectID   string
}

func NewClient(projectID string) *Client {
	ctx := context.Background()
	cPubSub, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	cSubscriber, err := pubsubapi.NewSubscriberClient(ctx)
	if err != nil {
		panic(err)
	}
	return &Client{
		cPubSub,
		cSubscriber,
		projectID,
	}
}

// メッセージを送信する
func (c *Client) Publish(
	ctx context.Context,
	topicID string,
	msg any,
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
