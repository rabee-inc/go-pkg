package rapi

import (
	"context"
	"net/http"
)

func NewHandlerMethod[I, O any](f func(ctx context.Context, param *I) (*O, error)) HandlerMethod[I] {
	return &handlerMethod[I, O]{
		ServiceFunc: f,
	}
}

type handlerMethod[I, O any] struct {
	InputFunc             func(ctx context.Context, r *http.Request, param any) error
	AfterInputFunc        func(ctx context.Context, r *http.Request, param *I) error
	ValidateFunc          func(ctx context.Context, param any) error
	AfterValidateFunc     func(ctx context.Context, param *I) error
	ServiceFunc           func(ctx context.Context, param *I) (*O, error)
	BeforeHandleErrorFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error
	HandleErrorFunc       func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error)
	RenderFunc            func(ctx context.Context, w http.ResponseWriter, r *http.Request, output any)
}

// --- handlerMethod implements ---

func (h *handlerMethod[I, O]) SetInputFunc(f func(ctx context.Context, r *http.Request, param any) error) {
	h.InputFunc = f
}

func (h *handlerMethod[I, O]) SetValidateFunc(f func(ctx context.Context, param any) error) {
	h.ValidateFunc = f
}

func (h *handlerMethod[I, O]) SetRenderFunc(f func(ctx context.Context, w http.ResponseWriter, r *http.Request, output any)) {
	h.RenderFunc = f
}

func (h *handlerMethod[I, O]) SetHandleErrorFunc(f func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error)) {
	h.HandleErrorFunc = f
}

func (h *handlerMethod[I, O]) BeforeHandleError(f func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error) {
	h.BeforeHandleErrorFunc = f
}

func (h *handlerMethod[I, O]) AfterInput(f func(ctx context.Context, r *http.Request, param *I) error) {
	h.AfterInputFunc = f
}

func (h *handlerMethod[I, O]) AfterValidate(f func(ctx context.Context, param *I) error) {
	h.AfterValidateFunc = f
}

// --- RouterElement implements ---

func (h *handlerMethod[I, O]) handleError(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	if h.BeforeHandleErrorFunc != nil {
		err = h.BeforeHandleErrorFunc(ctx, w, r, err)
	}
	if h.HandleErrorFunc != nil {
		h.HandleErrorFunc(ctx, w, r, err)
	}
}

func (h *handlerMethod[I, O]) GetHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var param I
		if h.InputFunc != nil {
			if err := h.InputFunc(ctx, r, &param); err != nil {
				h.handleError(ctx, w, r, err)
				return
			}
		}
		if h.AfterInputFunc != nil {
			if err := h.AfterInputFunc(ctx, r, &param); err != nil {
				h.handleError(ctx, w, r, err)
				return
			}
		}

		if h.ValidateFunc != nil {
			if err := h.ValidateFunc(ctx, &param); err != nil {
				h.handleError(ctx, w, r, err)
				return
			}
		}

		if h.AfterValidateFunc != nil {
			if err := h.AfterValidateFunc(ctx, &param); err != nil {
				h.handleError(ctx, w, r, err)
				return
			}
		}
		var output *O
		var err error
		if h.ServiceFunc != nil {
			output, err = h.ServiceFunc(ctx, &param)
			if err != nil {
				h.handleError(ctx, w, r, err)
				return
			}
		}

		if h.RenderFunc == nil {
			panic("RenderFunc is required")
		}
		h.RenderFunc(ctx, w, r, output)
	}
}

func (h *handlerMethod[I, O]) GetEmptyInput() any {
	return *new(I)
}

func (h *handlerMethod[I, O]) GetEmptyOutput() any {
	return *new(O)
}
