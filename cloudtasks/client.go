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
)

// ローカル実行時にタスクを直接叩く際のタイムアウト
const localTaskTimeout = 10 * time.Minute

type Client struct {
	client     *cloudtasks.Client
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
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return &Client{
		client:     client,
		port:       port,
		deploy:     deploy,
		projectID:  projectID,
		serviceID:  serviceID,
		locationID: locationID,
		authToken:  authToken,
	}
}

// クライアントを閉じる
func (c *Client) Close() error {
	return c.client.Close()
}

// リクエストをEnqueueする
func (c *Client) AddTask(ctx context.Context, queue string, path string, params any) error {
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
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": c.authToken,
		},
		Body: body,
	}
	return c.addTask(ctx, queue, req)
}

func (c *Client) addTask(ctx context.Context, queue string, aeReq *cloudtaskspb.AppEngineHttpRequest) error {
	if deploy.IsLocal() {
		return c.addTaskLocal(ctx, aeReq)
	}
	req := &cloudtaskspb.CreateTaskRequest{
		Parent: c.queuePath(queue),
		Task: &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: aeReq,
			},
		},
	}
	if _, err := c.client.CreateTask(ctx, req); err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// ローカルではCloud Tasksを経由せずに直接リクエストする
func (c *Client) addTaskLocal(ctx context.Context, aeReq *cloudtaskspb.AppEngineHttpRequest) error {
	url := fmt.Sprintf("http://localhost:%d%s", c.port, aeReq.GetRelativeUri())
	status, _, err := httpclient.PostBody(ctx, url, aeReq.GetBody(), &httpclient.HTTPOption{
		Headers: aeReq.GetHeaders(),
		Timeout: localTaskTimeout,
	})
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	if status != http.StatusOK {
		return log.Errore(ctx, "task http status: %d", status)
	}
	return nil
}

func (c *Client) queuePath(queue string) string {
	return fmt.Sprintf("projects/%s/locations/%s/queues/%s", c.projectID, c.locationID, queue)
}
