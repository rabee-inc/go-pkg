package maintenance

import (
	"net/http"

	"github.com/rabee-inc/go-pkg/renderer"
)

// Middleware ... Headerに関する機能を提供する
type Middleware struct {
	msg string
}

// Handle ... メンテナンスのハンドラ
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		renderer.Error(ctx, w, http.StatusServiceUnavailable, m.msg)
		return
	})
}

// NewMiddleware ... Middlewareを作成する
func NewMiddleware(msg string) *Middleware {
	if msg == "" {
		msg = "メンテナンス中のためご利用できません"
	}
	return &Middleware{
		msg: msg,
	}
}
