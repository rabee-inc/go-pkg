package parameter

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/util"
)

// GetURL ... リクエストからURLパラメータを取得する
func GetURL(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// GetURLByInt ... リクエストからURLパラメータをintで取得する
func GetURLByInt(ctx context.Context, r *http.Request, key string) (int, error) {
	str := chi.URLParam(r, key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Warningm(ctx, "strconv.Atoi", err)
		return num, err
	}
	return num, nil
}

// GetURLByInt64 ... リクエストからURLパラメータをint64で取得する
func GetURLByInt64(ctx context.Context, r *http.Request, key string) (int64, error) {
	str := chi.URLParam(r, key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Warningm(ctx, "strconv.ParseInt", err)
		return num, err
	}
	return num, nil
}

// GetURLByFloat64 ... リクエストからURLパラメータをfloat64で取得する
func GetURLByFloat64(ctx context.Context, r *http.Request, key string) (float64, error) {
	str := chi.URLParam(r, key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Warningm(ctx, "strconv.ParseFloat", err)
		return num, err
	}
	return num, nil
}

// GetForm ... リクエストからFormパラメータをstringで取得する
func GetForm(r *http.Request, key string) string {
	return r.FormValue(key)
}

// GetFormByInt ... リクエストからFormパラメータをintで取得する
func GetFormByInt(ctx context.Context, r *http.Request, key string) (int, error) {
	str := r.FormValue(key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		log.Warningm(ctx, "strconv.Atoi", err)
		return num, err
	}
	return num, nil
}

// GetFormByInt64 ... リクエストからFormパラメータをint64で取得する
func GetFormByInt64(ctx context.Context, r *http.Request, key string) (int64, error) {
	str := r.FormValue(key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Warningm(ctx, "strconv.ParseInt", err)
		return num, err
	}
	return num, nil
}

// GetFormByFloat64 ... リクエストからFormパラメータをfloat64で取得する
func GetFormByFloat64(ctx context.Context, r *http.Request, key string) (float64, error) {
	str := r.FormValue(key)
	if str == "" {
		return 0, nil
	}
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Warningm(ctx, "strconv.ParseFloat", err)
		return num, err
	}
	return num, nil
}

// GetFormByBool ... リクエストからFormパラメータをboolで取得する
func GetFormByBool(ctx context.Context, r *http.Request, key string) (bool, error) {
	str := r.FormValue(key)
	if str == "" {
		return false, nil
	}
	val, err := strconv.ParseBool(str)
	if err != nil {
		log.Warningm(ctx, "strconv.ParseInt", err)
		return val, err
	}
	return val, nil
}

// GetFormBySlice ... リクエストからFormパラメータをsliceで取得する
func GetFormBySlice(ctx context.Context, r *http.Request, key string) []string {
	sKey := fmt.Sprintf("%s[]", key)
	qs := r.URL.RawQuery
	vs := []string{}
	var err error
	for _, q := range strings.Split(qs, "&") {
		kv := strings.Split(q, "=")
		if len(kv) < 2 {
			continue
		}
		k := kv[0]
		k, err = url.QueryUnescape(k)
		if err != nil {
			log.Warningm(ctx, "url.QueryUnescape", err)
			continue
		}
		if k != sKey {
			continue
		}
		v := kv[1]
		v, err = url.QueryUnescape(v)
		if err != nil {
			log.Warningm(ctx, "url.QueryUnescape", err)
			continue
		}
		vs = append(vs, v)
	}
	return vs
}

// GetFormByIntSlice ... リクエストからFormパラメータをintのsliceで取得する
func GetFormByIntSlice(ctx context.Context, r *http.Request, key string) []int {
	strs := GetFormBySlice(ctx, r, key)
	nums := []int{}
	for _, str := range strs {
		num, err := strconv.Atoi(str)
		if err != nil {
			log.Warningm(ctx, "strconv.Atoi", err)
			continue
		}
		nums = append(nums, num)
	}
	return nums
}

// GetForms ... リクエストからFormパラメータを取得する
func GetForms(ctx context.Context, r *http.Request, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		err := log.Errore(ctx, "dst isn't a pointer")
		return err
	}

	paramType := reflect.TypeOf(dst).Elem()
	paramValue := reflect.ValueOf(dst).Elem()

	fieldCount := paramType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := paramType.Field(i)

		formTag := paramType.Field(i).Tag.Get("form")
		if util.IsZero(formTag) {
			continue
		}

		fieldValue := paramValue.FieldByName(field.Name)
		if !fieldValue.CanSet() {
			err := log.Warningc(ctx, http.StatusBadRequest, "fieldValue.CanSet")
			return err
		}
		switch field.Type.Kind() {
		case reflect.Int64, reflect.Int:
			val, err := GetFormByInt64(ctx, r, formTag)
			if err != nil {
				log.Warningm(ctx, "GetFormByInt64", err)
				return err
			}
			fieldValue.SetInt(val)
		case reflect.Float64:
			val, err := GetFormByFloat64(ctx, r, formTag)
			if err != nil {
				log.Warningm(ctx, "GetFormByFloat64", err)
				return err
			}
			fieldValue.SetFloat(val)
		case reflect.String:
			val := GetForm(r, formTag)
			fieldValue.SetString(val)
		case reflect.Bool:
			val, err := GetFormByBool(ctx, r, formTag)
			if err != nil {
				log.Warningm(ctx, "GetFormByBool", err)
				return err
			}
			fieldValue.SetBool(val)
		case reflect.Slice:
			switch {
			case field.Type == reflect.TypeOf([]string{}):
				val := GetFormBySlice(ctx, r, formTag)
				rv := reflect.ValueOf(val)
				fieldValue.Set(rv)
			case field.Type == reflect.TypeOf([]int{}):
				val := GetFormByIntSlice(ctx, r, formTag)
				rv := reflect.ValueOf(val)
				fieldValue.Set(rv)
			}
		}
	}
	return nil
}

// GetJSON ... リクエストからJSONパラメータを取得する
func GetJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)
	if err != nil {
		ctx := r.Context()
		log.Warningm(ctx, "dec.Decode", err)
		return err
	}
	return nil
}

// GetFormFile ... リクエストからファイルを取得する
func GetFormFile(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	return r.FormFile(key)
}
