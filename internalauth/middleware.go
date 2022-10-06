package internalauth

import (
	"fmt"
	"net/http"

	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/renderer"
)

type Middleware struct {
	Token string
}

func NewMiddleware(token string) *Middleware {
	return &Middleware{
		token,
	}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ah := r.Header.Get("Authorization")
		if ah == "" || ah != m.Token {
			msg := fmt.Sprintf("Internal auth error token: %s", ah)
			log.Warningf(ctx, msg)
			renderer.Error(ctx, w, http.StatusForbidden, msg)
			return
		}
		next.ServeHTTP(w, r)
	})
}
