package accesscontrol

import (
	"net/http"
	"strings"
)

type Middleware struct {
	origins []string
	header  string
}

func NewMiddleware(origins []string, headers []string) *Middleware {
	if len(origins) == 0 {
		origins = []string{}
	}
	if headers == nil {
		headers = []string{}
	}
	headers = append(headers, "Origin")
	headers = append(headers, "Content-Type")
	headers = append(headers, "Authorization")
	header := strings.Join(headers, ", ")
	return &Middleware{origins, header}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var host string
		if headers, ok := r.Header["Origin"]; ok {
			if len(headers) > 0 {
				host = headers[0]
			}
		}
		var origin string
		for _, o := range m.origins {
			if strings.Contains(o, host) {
				origin = o
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", m.header)
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) HandleWildcard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", m.header)
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
