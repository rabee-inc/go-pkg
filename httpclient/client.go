package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/rabee-inc/go-pkg/log"
)

const defaultTimeout time.Duration = 7 * time.Second

type HTTPOption struct {
	Headers map[string]string
	Timeout time.Duration
}

// Getリクエスト(URL)
func Get(ctx context.Context, url string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Getリクエスト(URL, param)
func GetForm(ctx context.Context, url string, param map[string]interface{}, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Warning(ctx, err)
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

// Getリクエスト(URL, QueryString)
func GetQueryString(ctx context.Context, url string, qs string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, url+"?"+qs, nil)
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Postリクエスト(URL, param)
func PostForm(ctx context.Context, url string, param map[string]interface{}, opt *HTTPOption) (int, []byte, error) {
	values := neturl.Values{}
	for key, value := range param {
		values.Add(key, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(values.Encode()))
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	return send(ctx, req, opt)
}

// Postリクエスト(URL, JSON)
func PostJSON(ctx context.Context, url string, param interface{}, res interface{}, opt *HTTPOption) (int, error) {
	jp, err := json.Marshal(param)
	if err != nil {
		log.Warning(ctx, err)
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jp))
	if err != nil {
		log.Warning(ctx, err)
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
	if len(body) > 0 {
		if status == http.StatusOK {
			perr := json.Unmarshal(body, res)
			if perr != nil {
				log.Warning(ctx, perr)
				err = perr
			}
		} else {
			errRes := map[string]interface{}{}
			err = json.Unmarshal(body, &errRes)
			log.Warningf(ctx, "%v", errRes)
		}
	}
	return status, err
}

// Postリクエスト(URL, Body)
func PostBody(ctx context.Context, url string, body []byte, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Putリクエスト(URL, JSON)
func PutJSON(ctx context.Context, url string, param interface{}, res interface{}, opt *HTTPOption) (int, error) {
	jp, err := json.Marshal(param)
	if err != nil {
		log.Warning(ctx, err)
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jp))
	if err != nil {
		log.Warning(ctx, err)
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
	if len(body) > 0 {
		if status == http.StatusOK {
			perr := json.Unmarshal(body, res)
			if perr != nil {
				log.Warning(ctx, perr)
				err = perr
			}
		} else {
			errRes := map[string]interface{}{}
			err = json.Unmarshal(body, &errRes)
			log.Warningf(ctx, "%v", errRes)
		}
	}
	return status, err
}

// Putリクエスト(URL, Body)
func PutBody(ctx context.Context, url string, body []byte, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Patchリクエスト(URL, JSON)
func PatchJSON(ctx context.Context, url string, param interface{}, res interface{}, opt *HTTPOption) (int, error) {
	jp, err := json.Marshal(param)
	if err != nil {
		log.Warning(ctx, err)
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jp))
	if err != nil {
		log.Warning(ctx, err)
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
	if len(body) > 0 {
		if status == http.StatusOK {
			perr := json.Unmarshal(body, res)
			if perr != nil {
				log.Warning(ctx, perr)
				err = perr
			}
		} else {
			errRes := map[string]interface{}{}
			err = json.Unmarshal(body, &errRes)
			log.Warningf(ctx, "%v", errRes)
		}
	}
	return status, err
}

// Patchリクエスト(URL, Body)
func PatchBody(ctx context.Context, url string, body []byte, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Deleteリクエスト(URL)
func Delete(ctx context.Context, url string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Deleteリクエスト(URL, param)
func DeleteForm(ctx context.Context, url string, param map[string]interface{}, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Warning(ctx, err)
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

// Deleteリクエスト(URL, QueryString)
func DeleteQueryString(ctx context.Context, url string, qs string, opt *HTTPOption) (int, []byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url+"?"+qs, nil)
	if err != nil {
		log.Warning(ctx, err)
		return 0, nil, err
	}

	if opt != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}
	return send(ctx, req, opt)
}

// Deleteリクエスト(URL, JSON)
func DeleteJSON(ctx context.Context, url string, param interface{}, res interface{}, opt *HTTPOption) (int, error) {
	jp, err := json.Marshal(param)
	if err != nil {
		log.Warning(ctx, err)
		return 0, err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jp))
	if err != nil {
		log.Warning(ctx, err)
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
	if len(body) > 0 {
		if status == http.StatusOK {
			perr := json.Unmarshal(body, res)
			if perr != nil {
				log.Warning(ctx, perr)
				err = perr
			}
		} else {
			errRes := map[string]interface{}{}
			err = json.Unmarshal(body, &errRes)
			log.Warningf(ctx, "%v", errRes)
		}
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
		log.Warning(ctx, err)
		return 0, nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Warning(ctx, err)
		return res.StatusCode, nil, nil
	}
	defer res.Body.Close()

	return res.StatusCode, body, nil
}
