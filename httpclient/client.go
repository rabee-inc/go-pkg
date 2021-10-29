package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/rabee-inc/go-pkg/log"
)

const defaultTimeout time.Duration = 20 * time.Second

// HTTPOption ... HTTP通信モジュールの追加設定
type HTTPOption struct {
	Headers map[string]string
	Timeout time.Duration
}

// Get ... Getリクエスト(URL)
func Get(ctx context.Context, url string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// GetForm ... Getリクエスト(URL, param)
func GetForm(ctx context.Context, url string, param map[string]interface{}, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	query := req.URL.Query()
	for key, value := range param {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = query.Encode()
	return send(ctx, req, opt)
}

// GetQueryString ... Getリクエスト(URL, QueryString)
func GetQueryString(ctx context.Context, url string, qs string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url+"?"+qs, nil)
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// PostForm ... Postリクエスト(URL, param)
func PostForm(ctx context.Context, url string, param map[string]interface{}, opt *HTTPOption) (int, []byte, error) {
	values := neturl.Values{}
	for key, value := range param {
		values.Add(key, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(values.Encode()))
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	return send(ctx, req, opt)
}

// PostJSON ... Postリクエスト(URL, JSON)
func PostJSON(ctx context.Context, url string, param interface{}, res interface{}, opt *HTTPOption) (int, error) {
	jp, err := json.Marshal(param)
	if err != nil {
		log.Warningm(ctx, "json.Marshal", err)
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jp))
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, err
	}

	if opt == nil {
		opt = &HTTPOption{
			Headers: map[string]string{},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range opt.Headers {
		req.Header.Set(key, value)
	}

	status, body, err := send(ctx, req, opt)
	if body != nil && len(body) > 0 {
		errRes := map[string]interface{}{}
		err = json.Unmarshal(body, &errRes)
		if err != nil {
			log.Warningm(ctx, "json.Unmarshal", err)
			return status, err
		}
		log.Warningf(ctx, "%v", errRes)
	}
	return status, err
}

// PostBody ... Postリクエスト(URL, Body)
func PostBody(ctx context.Context, url string, body []byte, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// PutJSON ... Putリクエスト(URL, JSON)
func PutJSON(ctx context.Context, url string, param interface{}, res interface{}, opt *HTTPOption) (int, error) {
	jp, err := json.Marshal(param)
	if err != nil {
		log.Warningm(ctx, "json.Marshal", err)
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jp))
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, err
	}

	if opt == nil {
		opt = &HTTPOption{
			Headers: map[string]string{},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range opt.Headers {
		req.Header.Set(key, value)
	}

	status, body, err := send(ctx, req, opt)
	if status != http.StatusOK {
		return status, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		log.Warningm(ctx, "json.Unmarshal", err)
		return status, err
	}
	return status, err
}

// PutBody ... Putリクエスト(URL, Body)
func PutBody(ctx context.Context, url string, body []byte, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Delete ... Deleteリクエスト(URL)
func Delete(ctx context.Context, url string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// DeleteForm ... Deleteリクエスト(URL, param)
func DeleteForm(ctx context.Context, url string, param map[string]interface{}, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	query := req.URL.Query()
	for key, value := range param {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	req.URL.RawQuery = query.Encode()
	return send(ctx, req, opt)
}

// DeleteQueryString ... Deleteリクエスト(URL, QueryString)
func DeleteQueryString(ctx context.Context, url string, qs string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url+"?"+qs, nil)
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// DeleteJSON ... Deleteリクエスト(URL, JSON)
func DeleteJSON(ctx context.Context, url string, param interface{}, res interface{}, opt *HTTPOption) (int, error) {
	jp, err := json.Marshal(param)
	if err != nil {
		log.Warningm(ctx, "json.Marshal", err)
		return 0, err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jp))
	if err != nil {
		log.Warningm(ctx, "http.NewRequest", err)
		return 0, err
	}

	if opt == nil {
		opt = &HTTPOption{
			Headers: map[string]string{},
		}
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range opt.Headers {
		req.Header.Set(key, value)
	}

	status, body, err := send(ctx, req, opt)
	if status != http.StatusOK {
		return status, err
	}

	err = json.Unmarshal(body, res)
	if err != nil {
		log.Warningm(ctx, "json.Unmarshal", err)
		return status, err
	}
	return status, err
}

func send(ctx context.Context, req *http.Request, opt *HTTPOption) (int, []byte, error) {
	client := http.Client{}
	if opt != nil && opt.Timeout > 0 {
		client.Timeout = opt.Timeout
	} else {
		client.Timeout = defaultTimeout
	}

	res, err := client.Do(req)
	if err != nil {
		log.Warningm(ctx, "client.Do", err)
		return 0, nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Warningm(ctx, "ioutil.ReadAll", err)
		return res.StatusCode, nil, nil
	}
	defer res.Body.Close()

	return res.StatusCode, body, nil
}
