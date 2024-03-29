package jsonrpc2

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rabee-inc/go-pkg/httpclient"
	"github.com/rabee-inc/go-pkg/log"
)

type Client struct {
	URL      string
	Headers  map[string]string
	Requests []*ClientRequest
}

func NewClient(url string, headers map[string]string) *Client {
	return &Client{
		url,
		headers,
		[]*ClientRequest{},
	}
}

// JSONRPC2のリクエストを登録する
func (c *Client) AddRequest(ctx context.Context, id string, method string, params any) error {
	rawParams, err := c.marshalRawMessage(ctx, params)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	c.Requests = append(c.Requests, &ClientRequest{
		Version: version,
		ID:      id,
		Method:  method,
		Params:  rawParams,
	})
	return nil
}

// JSONRPC2のシングルリクエストを行う
func (c *Client) DoSingle(ctx context.Context, method string, params any) (*json.RawMessage, *ErrorResponse, error) {
	rawParams, err := c.marshalRawMessage(ctx, params)
	if err != nil {
		log.Error(ctx, err)
		return nil, nil, err
	}

	req := &ClientRequest{
		Version: version,
		ID:      "single",
		Method:  method,
		Params:  rawParams,
	}
	var res ClientResponse
	status, err := httpclient.PostJSON(ctx, c.URL, req, &res, &httpclient.HTTPOption{Headers: c.Headers})
	if err != nil {
		log.Error(ctx, err)
		return nil, nil, err
	}
	if status != http.StatusOK {
		err := log.Errore(ctx, "httpclient.PostJSON status: %d", status)
		return nil, nil, err
	}
	return res.Result, res.Error, nil
}

// JSONRPC2のバッチリクエストを行う
func (c *Client) DoBatch(ctx context.Context) ([]*ClientResponse, error) {
	var res []*ClientResponse
	status, err := httpclient.PostJSON(ctx, c.URL, c.Requests, &res, nil)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	if status != http.StatusOK {
		err := log.Errore(ctx, "httpclient.PostJSON status: %d", status)
		return nil, err
	}
	return res, nil
}

func (c *Client) marshalRawMessage(ctx context.Context, params any) (*json.RawMessage, error) {
	bParams, err := json.Marshal(params)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	rawMessage := json.RawMessage(bParams)
	return &rawMessage, nil
}
