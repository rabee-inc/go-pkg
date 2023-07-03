package rapi

import (
	"net/http"
	"reflect"

	"github.com/go-chi/chi"
	"github.com/rabee-inc/go-pkg/util"
)

func FillURLParam(r *http.Request, param any) {
	_ = util.EachTaggedFields(param, "url", func(tagValue string, reflectParam reflect.Value, fieldNum int) error {
		urlParam := chi.URLParam(r, tagValue)
		reflectParam.Field(fieldNum).SetString(urlParam)
		return nil
	})
}
