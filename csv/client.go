package csv

import (
	"context"
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/jszwec/csvutil"
	"github.com/rabee-inc/go-pkg/bytesutil"
	"github.com/rabee-inc/go-pkg/httpclient"
	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/stringutil"
)

// レスポンス形式に変換する
func ToResponse(ctx context.Context, srcs interface{}) ([][]string, error) {
	bytes, err := csvutil.Marshal(srcs)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	str := bytesutil.ToStr(bytes)

	r := csv.NewReader(strings.NewReader(str))
	rows, err := r.ReadAll()
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return rows, nil
}

// URLからCSVデータを取得する
func GetByURL(ctx context.Context, url string, dsts interface{}) error {
	status, body, err := httpclient.Get(ctx, url, nil)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	if status != http.StatusOK {
		err = log.Errore(ctx, "get csv request error: %d", status)
		return err
	}
	err = GetByBytes(ctx, body, dsts)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// 文字列からCSVデータを取得する
func GetByStr(ctx context.Context, str string, dsts interface{}) error {
	bytes := stringutil.ToBytes(str)
	err := GetByBytes(ctx, bytes, dsts)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}

// バイト列からCSVデータを取得する
func GetByBytes(ctx context.Context, bytes []byte, dsts interface{}) error {
	err := csvutil.Unmarshal(bytes, dsts)
	if err != nil {
		log.Error(ctx, err)
		return err
	}
	return nil
}
