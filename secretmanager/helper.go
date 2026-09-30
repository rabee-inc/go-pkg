package secretmanager

import (
	"context"
	"fmt"
	"os"

	sm "cloud.google.com/go/secretmanager/apiv1"
	sm_pb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	"github.com/rabee-inc/go-pkg/bytesutil"
)

// LoadSecret ... Secret Manager の値を環境変数に読み込む
func LoadSecret(projectID string, params []*LoadSecretParam) {
	ctx := context.Background()
	cSecretManager, err := sm.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer cSecretManager.Close()

	for _, param := range params {
		version := param.Version
		if version == "" {
			version = "latest"
		}
		name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, param.Key, version)
		request := &sm_pb.AccessSecretVersionRequest{
			Name: name,
		}
		result, err := cSecretManager.AccessSecretVersion(ctx, request)
		if err != nil {
			panic(err)
		}
		v := bytesutil.ToStr(result.GetPayload().GetData())
		if err := os.Setenv(param.Key, v); err != nil {
			panic(err)
		}
	}
}
