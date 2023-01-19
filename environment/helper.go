package environment

import (
	"fmt"

	"os"

	"github.com/rabee-inc/go-pkg/deploy"
	"gopkg.in/yaml.v3"
)

// 環境変数を読み込む
func Load(envFilePath string) {
	// 環境変数設定ファイル読み込み
	file, err := os.ReadFile(envFilePath)
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
