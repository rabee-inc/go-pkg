package jwtauth

import (
	"context"
	"net/http"

	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/renderer"
)

// Middleware ... JWT認証のミドルウェア
type Middleware struct {
	Svc Service
}

// Handle ... JWT認証をする
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Headerを取得
		ah := r.Header.Get("Authorization")
		if ah == "" {
			err := log.Warningc(ctx, http.StatusForbidden, "no Authorization header")
			m.renderError(ctx, w, err)
			return
		}

		// 認証
		userID, claims, err := m.Svc.Authentication(ctx, ah)
		if err != nil {
			log.Warning(ctx, err)
			m.renderError(ctx, w, err)
			return
		}

		// 認証結果を設定
		ctx = setUserID(ctx, userID)
		log.Debugf(ctx, "UserID: %s", userID)

		ctx = setClaims(ctx, claims)
		log.Debugf(ctx, "Claims: %v", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) renderError(ctx context.Context, w http.ResponseWriter, err error) {
	var status int
	if cd, ok := errcode.Get(err); ok {
		status = cd
	} else {
		status = http.StatusForbidden
	}
	renderer.Error(ctx, w, status, err.Error())
}

// NewMiddleware ... Middlewareを作成する
func NewMiddleware(svc Service) *Middleware {
	return &Middleware{
		Svc: svc,
	}
}
