package images

import (
	"context"

	"github.com/rabee-inc/go-pkg/cloudpubsub"
	"github.com/rabee-inc/go-pkg/log"
)

type Client struct {
	cPubSub          *cloudpubsub.Client
	converterTopicID string
	generatorTopicID string
}

func NewClient(
	cPubSub *cloudpubsub.Client,
	dstEndpoint string,
) *Client {
	return &Client{
		cPubSub,
		ConverterTopicID,
		GeneratorTopicID,
	}
}

func NewClientWithOption(
	cPubSub *cloudpubsub.Client,
	reqOption *ClientOption,
) *Client {
	option := &ClientOption{
		ConverterTopicID: ConverterTopicID,
		GeneratorTopicID: GeneratorTopicID,
	}
	if reqOption != nil && reqOption.ConverterTopicID != "" {
		option.ConverterTopicID = reqOption.ConverterTopicID
	}
	if reqOption != nil && reqOption.GeneratorTopicID != "" {
		option.GeneratorTopicID = reqOption.GeneratorTopicID
	}
	return &Client{
		cPubSub,
		option.ConverterTopicID,
		option.GeneratorTopicID,
	}
}

// 画像変換リクエストを送信する
func (c *Client) SendConvertRequest(
	ctx context.Context,
	key string,
	sourceID string,
	sources []*Object,
	dstFilePath string,
	dstEndpoint string,
	dstEndpointAuthToken string,
) error {
	srcURLs := []string{}
	for _, source := range sources {
		if source == nil || source.URL == "" {
			continue
		}
		srcURLs = append(srcURLs, source.URL)
	}
	if len(srcURLs) == 0 {
		return nil
	}
	src := &ConvertRequest{
		Key:                  key,
		SourceID:             sourceID,
		SourceURLs:           srcURLs,
		DstFilePath:          dstFilePath,
		DstEndpoint:          dstEndpoint,
		DstEndpointAuthToken: dstEndpointAuthToken,
	}
	err := c.cPubSub.Publish(ctx, c.converterTopicID, src)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// 画像作成リクエストを送信する
func (c *Client) SendGenerateRequest(
	ctx context.Context,
	key string,
	sourceID string,
	sourceURL string,
	width int,
	height int,
	dstFilePath string,
) error {
	if sourceID == "" || sourceURL == "" || dstFilePath == "" {
		err := log.Errore(ctx, "invalid param sourceID: %s, sourceURL: %s, dstFilePath: %s", sourceID, sourceURL, dstFilePath)
		return err
	}
	src := &GenerateRequest{
		Key:         key,
		SourceID:    sourceID,
		SourceURL:   sourceURL,
		Width:       width,
		Height:      height,
		DstFilePath: dstFilePath,
	}
	err := c.cPubSub.Publish(ctx, c.generatorTopicID, src)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}
