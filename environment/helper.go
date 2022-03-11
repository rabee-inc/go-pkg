package environment

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"gopkg.in/yaml.v3"

	"github.com/rabee-inc/go-pkg/bytesutil"
	"github.com/rabee-inc/go-pkg/deploy"
)

// Load ... 環境変数を読み込む
func Load(envFilePath string) {
	// 環境変数設定ファイル読み込み
	file, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		panic(err)
	}
	val := &Variable{}
	err = yaml.Unmarshal(file, &val)
	if err != nil {
		panic(err)
	}

	// 値を設定
	var src map[string]string
	if deploy.IsLocal() {
		src = val.Local
		src["DEPLOY"] = "local"
	} else if deploy.IsStaging() {
		src = val.Staging
	} else if deploy.IsProduction() {
		src = val.Production
	} else {
		panic(fmt.Errorf("invalid deploy: %s", os.Getenv("DEPLOY")))
	}

	for k, v := range src {
		err = os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}
}

func LoadSecret(projectID string, params []*LoadSecretParam) {
	ctx := context.Background()
	cSecretManager, err := secretmanager.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	for _, param := range params {
		if param.Version == "" {
			param.Version = "latest"
		}
		name := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, param.Key, param.Version)
		request := &secretmanagerpb.AccessSecretVersionRequest{
			Name: name,
		}
		result, err := cSecretManager.AccessSecretVersion(ctx, request)
		if err != nil {
			panic(err)
		}
		v := bytesutil.ToStr(result.Payload.GetData())
		err = os.Setenv(param.Key, v)
		if err != nil {
			panic(err)
		}
	}
}
