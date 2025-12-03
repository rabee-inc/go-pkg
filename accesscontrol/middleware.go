package accesscontrol

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Middleware struct {
	originRegexps []*regexp.Regexp
	header        string
}

func NewMiddleware(origins []string, headers []string) *Middleware {
	originRegexps := []*regexp.Regexp{}
	for _, origin := range origins {
		origin = strings.ReplaceAll(origin, ".", "\\.")
		origin = strings.ReplaceAll(origin, "*", ".*")
		pattern := fmt.Sprintf("^%s$", origin)
		originRegexp := regexp.MustCompile(pattern)
		originRegexps = append(originRegexps, originRegexp)
	}
	if headers == nil {
		headers = []string{}
	}
	headers = append(headers, "Origin")
	headers = append(headers, "Content-Type")
	headers = append(headers, "Authorization")
	header := strings.Join(headers, ", ")
	return &Middleware{originRegexps, header}
}

func (m *Middleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := m.GetOriginValue(r.Header.Get("Origin"))
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

func (m *Middleware) GetOriginValue(requestOrigin string) string {
	var origin string
	if len(m.originRegexps) == 0 {
		// 何も指定されていなかったらワイルドカードを返す
		origin = "*"
	} else {
		for _, originRegexp := range m.originRegexps {
			if originRegexp.MatchString(requestOrigin) {
				origin = requestOrigin
				break
			}
		}
	}
	return origin
}
