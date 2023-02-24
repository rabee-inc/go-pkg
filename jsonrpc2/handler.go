package jsonrpc2

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/log"
	"github.com/unrolled/render"
)

type Handler struct {
	actions map[string]Action
}

func NewHandler() *Handler {
	return &Handler{
		map[string]Action{},
	}
}

// JSONRPC2のリクエストを登録する
func (h *Handler) Register(method string, action Action) {
	if method == "" || action == nil {
		panic(fmt.Errorf("invalid method name: %s, action: %v", method, action))
	}
	h.actions[method] = action
}

// JSONRPC2のリクエストをハンドルする
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", contentType)

	// POSTで送信されていること
	if r.Method != "POST" {
		log.SetResponseStatus(ctx, http.StatusMethodNotAllowed)
		h.renderError(ctx, w, http.StatusMethodNotAllowed, "invalid http method: %s", r.Method)
		return
	}

	// リクエストのContent-TypeもしくはAcceptがapplication/jsonであること
	ct := r.Header.Get("Content-Type")
	if ct != contentType {
		log.SetResponseStatus(ctx, http.StatusUnsupportedMediaType)
		h.renderError(ctx, w, http.StatusUnsupportedMediaType, "invalid http header content-type: %s", ct)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.SetResponseStatus(ctx, http.StatusBadRequest)
		h.renderError(ctx, w, http.StatusBadRequest, "read http body error: %s", err.Error())
		return
	}

	err = h.handleSingleRequest(ctx, w, r, data)
	if err != nil {
		err = h.handleBatchRequest(ctx, w, r, data)
	}
	if err != nil {
		log.SetResponseStatus(ctx, http.StatusBadRequest)
		h.renderError(ctx, w, http.StatusBadRequest, "parse json error: %s", err.Error())
		return
	}
	log.SetResponseStatus(ctx, http.StatusOK)
}

func (h *Handler) handleSingleRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, data []byte) error {
	var req request
	err := json.Unmarshal(data, &req)
	if err != nil {
		return err
	}
	res := h.handleRequest(ctx, r, req)
	encoder := json.NewEncoder(w)
	return encoder.Encode(res)
}

func (h *Handler) handleBatchRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, data []byte) error {
	var reqs []request
	err := json.Unmarshal(data, &reqs)
	if err != nil {
		return err
	}

	var responses []response
	ch := make(chan response, len(reqs))
	for _, req := range reqs {
		go func(_req request) {
			ch <- h.handleRequest(ctx, r, _req)
		}(req)
	}
	for i := 0; i < len(reqs); i++ {
		responses = append(responses, <-ch)
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(responses)
}

func (h *Handler) handleRequest(ctx context.Context, r *http.Request, req request) response {
	if !req.isValid() {
		return h.renderErrorJSON(ctx, req.ID, http.StatusBadRequest, "invalid jsonrpc2 params: %s", spew.Sdump(req))
	}

	action := h.actions[req.Method]
	if action == nil {
		return h.renderErrorJSON(ctx, req.ID, http.StatusBadRequest, "method not found: %s", req.Method)
	}

	params, err := action.DecodeParams(ctx, req.Params)
	if err != nil {
		return h.renderErrorJSON(ctx, req.ID, http.StatusBadRequest, "invalid params: %s", err.Error())
	}

	result, err := action.Exec(ctx, req.Method, params)
	if err != nil {
		code, ok := errcode.Get(err)
		if !ok {
			code = http.StatusInternalServerError
		}
		return h.renderErrorJSON(ctx, req.ID, code, err.Error())
	}

	return newResponse(req.ID, result)
}

func (h *Handler) renderError(ctx context.Context, w http.ResponseWriter, status int, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	switch status {
	case http.StatusBadRequest:
		log.Warningf(ctx, msg)
	case http.StatusForbidden:
		log.Warningf(ctx, msg)
	case http.StatusNotFound:
		log.Warningf(ctx, msg)
	default:
		log.Errorf(ctx, msg)
	}
	render.New().Text(w, status, msg)
}

func (h *Handler) renderErrorJSON(ctx context.Context, rpcID string, rpcStatus int, format string, a ...any) response {
	msg := fmt.Sprintf(format, a...)
	switch rpcStatus {
	case http.StatusBadRequest:
		log.Warningf(ctx, msg)
	case http.StatusForbidden:
		log.Warningf(ctx, msg)
	case http.StatusNotFound:
		log.Warningf(ctx, msg)
	default:
		log.Errorf(ctx, msg)
	}
	return newErrorResponse(rpcID, rpcStatus, msg)
}
