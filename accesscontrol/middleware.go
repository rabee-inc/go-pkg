package accesscontrol

import (
	"net/http"
	"strings"
)

// Middleware ... Headerに関する機能を提供する
type Middleware struct {
	hdsStr string
}

// Handle ... CORS対応
func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", m.hdsStr)
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NewMiddleware ... Middlewareを作成する
func NewMiddleware(headers []string) *Middleware {
	hds := []string{
		"Origin",
		"Content-Type",
		"Authorization",
	}
	for _, header := range headers {
		hds = append(hds, header)
	}
	hdsStr := strings.Join(hds, ", ")
	return &Middleware{
		hdsStr: hdsStr,
	}
}
