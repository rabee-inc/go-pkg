package log

import (
	"net/http"

	"github.com/rabee-inc/go-pkg/stringutil"
	"github.com/rabee-inc/go-pkg/timeutil"
)

// Middleware ... ロガー
type Middleware struct {
	Writer         Writer
	MinOutSeverity Severity
}

// Handle ... ロガーを初期化する
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startAt := timeutil.Now()

		// ロガーをContextに設定
		traceID := stringutil.UniqueID()
		logger := NewLogger(m.Writer, m.MinOutSeverity, traceID)
		ctx := r.Context()
		ctx = SetLogger(ctx, logger)

		// Panicのハンドリングを設定
		defer func() {
			if rcvr := recover(); rcvr != nil {
				msg := Panic(ctx, rcvr)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(msg))

				// 実行時間を計算
				endAt := timeutil.Now()
				dr := endAt.Sub(startAt)

				// リクエストログを出力
				logger.WriteRequest(r, endAt, dr)
			}
		}()

		// 実行
		next.ServeHTTP(w, r.WithContext(ctx))

		// 実行時間を計算
		endAt := timeutil.Now()
		dr := endAt.Sub(startAt)

		// リクエストログを出力
		logger.WriteRequest(r, endAt, dr)
	})
}

// NewMiddleware ... ミドルウェアを作成する
func NewMiddleware(writer Writer, minOutSeverity string) *Middleware {
	mos := NewSeverity(minOutSeverity)
	return &Middleware{
		Writer:         writer,
		MinOutSeverity: mos,
	}
}
