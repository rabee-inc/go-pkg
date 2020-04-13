package httpheader

import (
	"context"
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/rabee-inc/go-pkg/log"
)

const (
	headerKeySample string = "X-Sample"
)

type service struct {
}

func (s *service) Get(ctx context.Context, r *http.Request) (Params, error) {
	h := Params{
		// EDIT: ここに任意のヘッダーを入れる
		Sample: r.Header.Get(headerKeySample),
	}

	v := validator.New()
	if err := v.Struct(h); err != nil {
		log.Warningf(ctx, "Header param validation error: %s", err.Error())
		return h, err
	}

	return h, nil
}

// NewService ... Serviceを作成する
func NewService() Service {
	return &service{}
}
