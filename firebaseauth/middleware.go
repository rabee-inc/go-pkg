package firebaseauth

import (
	"net/http"

	"github.com/rabee-inc/go-pkg/log"
	"github.com/rabee-inc/go-pkg/renderer"
)

type Middleware struct {
	sFirebaseAuth Service
	optional      bool
}

func NewMiddleware(sFirebaseAuth Service, optional bool) *Middleware {
	return &Middleware{
		sFirebaseAuth,
		optional,
	}
}

// Handle ... Firebase認証をする
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Headerを取得
		ah := r.Header.Get("Authorization")
		if ah == "" {
			if !m.optional {
				log.Warningf(ctx, "no authorization header")
				renderer.Error(ctx, w, http.StatusUnauthorized, "アカウントの認証に失敗しました")
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = setAuthHeader(ctx, ah)

		// 認証
		userID, claims, err := m.sFirebaseAuth.Authentication(ctx, ah)
		if err != nil {
			log.Warning(ctx, err)
			renderer.Error(ctx, w, http.StatusUnauthorized, "アカウントの認証に失敗しました")
			return
		}

		// 認証結果を設定
		ctx = setUserID(ctx, userID)
		log.Infof(ctx, "UserID: %s", userID)

		ctx = setClaims(ctx, claims)
		log.Debugf(ctx, "Claims: %v", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
