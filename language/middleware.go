package language

import (
	"net/http"

	"github.com/rabee-inc/go-pkg/log"
)

// Middleware ... 言語のミドルウェア
type Middleware struct {
	headerKey  string
	defaultKey Key
}

// Handle ... 言語を設定する
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Headerを取得
		h := r.Header.Get(m.headerKey)
		log.Debugf(ctx, "%s: %s", m.headerKey, h)
		h = trim(h)

		var key Key
		if h == "" {
			key = m.defaultKey
		} else {
			key = Key(h)
		}

		// Contextに設定
		ctx = setKey(ctx, key)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewMiddleware ... Middlewareを作成する
func NewMiddleware(headerKey string, defaultKey Key) *Middleware {
	return &Middleware{
		headerKey:  headerKey,
		defaultKey: defaultKey,
	}
}
