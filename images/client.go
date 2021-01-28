package images

import (
	"context"

	"github.com/rabee-inc/go-pkg/cloudpubsub"
	"github.com/rabee-inc/go-pkg/log"
)

// Client ... クライアント
type Client struct {
	psCli *cloudpubsub.Client
}

// SendConvertRequest ... 画像変換リクエストを送信する
func (c *Client) SendConvertRequest(
	ctx context.Context,
	key string,
	sourceID string,
	sources []*Object,
	dstFilePath string) error {
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
		Key:         key,
		SourceID:    sourceID,
		SourceURLs:  srcURLs,
		DstFilePath: dstFilePath,
	}
	err := c.psCli.Publish(ctx, ConverterTopicID, src)
	if err != nil {
		log.Errorm(ctx, "c.psCli.Publish", err)
		return err
	}
	return nil
}

// SendGenerateRequest ... 画像作成リクエストを送信する
func (c *Client) SendGenerateRequest(
	ctx context.Context,
	key string,
	sourceID string,
	sourceURL string,
	width int,
	height int,
	dstFilePath string) error {
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
	err := c.psCli.Publish(ctx, GeneratorTopicID, src)
	if err != nil {
		log.Errorm(ctx, "c.psCli.Publish", err)
		return err
	}
	return nil
}

// NewClient ... クライアントを作成する
func NewClient(psCli *cloudpubsub.Client) *Client {
	return &Client{
		psCli: psCli,
	}
}
