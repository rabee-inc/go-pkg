package bigquery

import (
	"context"
	"reflect"

	"cloud.google.com/go/bigquery"
	"github.com/rabee-inc/go-pkg/log"
	"google.golang.org/api/iterator"
)

type Client struct {
	client *bigquery.Client
}

func NewClient(projectID string) *Client {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return &Client{client}
}

// クエリを実行し、データを取得する
func (c *Client) List(ctx context.Context, query string, limit int, cursor string, dsts any) (string, error) {
	q := c.client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		log.Error(ctx, err)
		return "", err
	}
	if pageInfo := it.PageInfo(); pageInfo != nil {
		pageInfo.MaxSize = limit
		pageInfo.Token = cursor
	}

	rv := reflect.Indirect(reflect.ValueOf(dsts))
	rrt := rv.Type().Elem().Elem()
	i := 0
	for {
		i++
		v := reflect.New(rrt).Interface()
		err = it.Next(v)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Error(ctx, err)
			return "", err
		}
		rrv := reflect.ValueOf(v)
		rv.Set(reflect.Append(rv, rrv))
		if i == limit {
			break
		}
	}
	var token string
	if pageInfo := it.PageInfo(); pageInfo != nil {
		token = pageInfo.Token
	}
	return token, nil
}
