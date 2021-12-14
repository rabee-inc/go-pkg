package cloudpubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pubsub "cloud.google.com/go/pubsub"
	pubsubapi "cloud.google.com/go/pubsub/apiv1"
	"google.golang.org/api/option"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/rabee-inc/go-pkg/log"
)

// Client ... Pub/Subのクライアント
type Client struct {
	projectID string
	psClient  *pubsub.Client
	psaClient *pubsubapi.SubscriberClient
	topics    map[string]*pubsub.Topic
}

// Publish ... Pushを利用してメッセージを送信する
func (c *Client) Publish(ctx context.Context, topicID string, src interface{}) error {
	topic, ok := c.topics[topicID]
	if !ok {
		return log.Errore(ctx, "no generated topic: %s", topicID)
	}
	bSrc, err := json.Marshal(src)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	result := topic.Publish(ctx, &pubsub.Message{
		Data: bSrc,
	})
	if result == nil {
		err = log.Errore(ctx, "not return publishResult for topic: %s", topicID)
		return err
	}
	_, err = result.Get(ctx)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// Receive ... Pullを利用してメッセージを受信する
func (c *Client) Receive(ctx context.Context, subID string, maxMessageCount int, dsts interface{}) error {
	ackIDs, msgs, err := c.sendPull(ctx, subID, maxMessageCount)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	if len(ackIDs) > 0 {
		err = c.sendAcknowledge(ctx, subID, ackIDs)
		if err != nil {
			log.Error(ctx, err)
			return err
		}
	}
	bMsg, err := json.Marshal(msgs)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	err = json.Unmarshal(bMsg, dsts)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

func (c *Client) sendPull(ctx context.Context, subID string, maxMessageCount int) ([]string, []json.RawMessage, error) {
	res, err := c.psaClient.Pull(ctx, &pubsubpb.PullRequest{
		Subscription:      c.generateSub(subID),
		ReturnImmediately: true,
		MaxMessages:       int32(maxMessageCount),
	})
	if err != nil {
		log.Error(ctx, err)
		return nil, nil, err
	}
	if len(res.ReceivedMessages) == 0 {
		return []string{}, []json.RawMessage{}, nil
	}
	ackIDs := []string{}
	dsts := []json.RawMessage{}
	for _, receivedMessage := range res.ReceivedMessages {
		ackIDs = append(ackIDs, receivedMessage.AckId)
		if message := receivedMessage.Message; message != nil {
			dsts = append(dsts, message.Data)
		}
	}
	return ackIDs, dsts, nil
}

func (c *Client) sendAcknowledge(ctx context.Context, subID string, ackIDs []string) error {
	err := c.psaClient.Acknowledge(ctx, &pubsubpb.AcknowledgeRequest{
		Subscription: c.generateSub(subID),
		AckIds:       ackIDs,
	})
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

func (c *Client) generateSub(subID string) string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", c.projectID, subID)
}

// NewClient ... Pub/Subのクライアントを取得する
func NewClient(projectID string, topicIDs []string) *Client {
	// Clientを作成
	ctx := context.Background()
	gOpt := option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                1 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}))
	psClient, err := pubsub.NewClient(ctx, projectID, gOpt)
	if err != nil {
		panic(err)
	}
	psaClient, err := pubsubapi.NewSubscriberClient(ctx, gOpt)
	if err != nil {
		panic(err)
	}

	// Topicsを作成
	topics := map[string]*pubsub.Topic{}
	for _, topicID := range topicIDs {
		topics[topicID] = psClient.Topic(topicID)
	}
	return &Client{
		projectID: projectID,
		psClient:  psClient,
		psaClient: psaClient,
		topics:    topics,
	}
}
