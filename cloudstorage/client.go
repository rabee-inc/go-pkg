package cloudstorage

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/timeutil"
	"github.com/vincent-petithory/dataurl"
)

type Client struct {
	client       *storage.Client
	bucketHandle *storage.BucketHandle
	bucket       string
}

func NewClient(
	bucket string,
) *Client {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	bucketHandle := client.Bucket(bucket)
	return &Client{
		client,
		bucketHandle,
		bucket,
	}
}

// DataURLのファイルをアップロードする
func (c *Client) UploadForDataURL(
	ctx context.Context,
	path string,
	name string,
	cacheMode *CacheMode,
	dataURL string,
) (string, error) {
	// Base64をデコード
	res, err := dataurl.DecodeString(dataURL)
	if err != nil {
		log.Warning(ctx, err)
		err = errcode.Set(err, http.StatusBadRequest)
		return "", err
	}

	// アップロード
	return c.Upload(ctx, path, name, res.ContentType(), cacheMode, res.Data)
}

// ファイルをアップロードする
func (c *Client) Upload(
	ctx context.Context,
	path string,
	name string,
	contentType string,
	cacheMode *CacheMode,
	data []byte,
) (string, error) {
	// Writerを作成
	w := c.client.
		Bucket(c.bucket).
		Object(strings.Join([]string{path, name}, "/")).
		NewWriter(ctx)

	// ContentTypeを設定
	w.ContentType = contentType

	// Cache-Controllを設定
	if cacheMode != nil {
		var cc string
		if cacheMode.Disabled {
			cc = "no-cache"
		} else {
			cc = fmt.Sprintf("max-age=%d", cacheMode.Expire/time.Second)
		}
		w.CacheControl = cc
	}
	w.ChunkSize = ChunkSize

	// アップロード
	if _, err := w.Write(data); err != nil {
		log.Error(ctx, err)
		return "", err
	}
	if err := w.Close(); err != nil {
		log.Error(ctx, err)
		return "", err
	}

	// URLを作成
	url := GenerateFileURL(c.bucket, path, name)
	return url, nil
}

// 指定ファイルのReaderを取得する
func (c *Client) GetReader(
	ctx context.Context,
	path string,
) (*storage.Reader, error) {
	reader, err := c.client.
		Bucket(c.bucket).
		Object(path).
		NewReader(ctx)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return reader, nil
}

// バケット名
func (c *Client) GetBucket() string {
	return c.bucket
}

func (c *Client) GetDownloadSignedURL(
	ctx context.Context,
	path string,
	expire time.Duration,
) (string, error) {
	expires := timeutil.Now().Add(expire)
	opts := &storage.SignedURLOptions{
		Expires: expires,
	}
	opts.Method = http.MethodGet
	singedURL, err := c.bucketHandle.SignedURL(path, opts)
	if err != nil {
		log.Error(ctx, err)
		return "", err
	}
	return singedURL, nil
}

func (c *Client) GetUploadSignedURL(
	ctx context.Context,
	path string,
	contentType string,
	expire time.Duration,
) (string, error) {
	expires := timeutil.Now().Add(expire)
	opts := &storage.SignedURLOptions{
		Expires: expires,
	}
	opts.Method = http.MethodPut
	opts.ContentType = contentType
	singedURL, err := c.bucketHandle.SignedURL(path, opts)
	if err != nil {
		log.Error(ctx, err)
		return "", err
	}
	return singedURL, nil
}
