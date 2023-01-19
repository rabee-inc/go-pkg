package cloudtasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/rabee-inc/go-pkg/deploy"
	"github.com/rabee-inc/go-pkg/httpclient"
	"github.com/rabee-inc/go-pkg/log"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Client struct {
	cli        *cloudtasks.Client
	port       int
	deploy     string
	projectID  string
	serviceID  string
	locationID string
	authToken  string
}

func NewClient(
	port int,
	deploy string,
	projectID string,
	serviceID string,
	locationID string,
	authToken string) *Client {
	ctx := context.Background()
	gOpt := option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                1 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}))
	cli, err := cloudtasks.NewClient(ctx, gOpt)
	if err != nil {
		panic(err)
	}
	return &Client{
		cli,
		port,
		deploy,
		projectID,
		serviceID,
		locationID,
		authToken,
	}
}

// リクエストをEnqueueする
func (c *Client) AddTask(ctx context.Context, queue string, path string, params interface{}) error {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": c.authToken,
	}
	body, err := json.Marshal(params)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	req := &cloudtaskspb.AppEngineHttpRequest{
		AppEngineRouting: &cloudtaskspb.AppEngineRouting{
			Service: c.serviceID,
		},
		HttpMethod:  cloudtaskspb.HttpMethod_POST,
		RelativeUri: path,
		Headers:     headers,
		Body:        body,
	}
	return c.addTask(ctx, queue, req)
}

func (c *Client) addTask(ctx context.Context, queue string, aeReq *cloudtaskspb.AppEngineHttpRequest) error {
	if deploy.IsLocal() {
		url := fmt.Sprintf("http://localhost:%d%s", c.port, aeReq.RelativeUri)
		status, _, err := httpclient.PostBody(ctx, url, aeReq.Body, &httpclient.HTTPOption{
			Headers: aeReq.Headers,
		})
		if err != nil {
			log.Error(ctx, err)
			return err
		}
		if status != http.StatusOK {
			err = log.Errore(ctx, "task http status: %d", status)
			return err
		}
	} else {
		req := &cloudtaskspb.CreateTaskRequest{
			Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", c.projectID, c.locationID, queue),
			Task: &cloudtaskspb.Task{
				MessageType: &cloudtaskspb.Task_AppEngineHttpRequest{
					AppEngineHttpRequest: aeReq,
				},
			},
		}
		_, err := c.cli.CreateTask(ctx, req)
		if err != nil {
			log.Error(ctx, err)
			return err
		}
	}
	return nil
}
