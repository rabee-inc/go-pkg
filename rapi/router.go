package rapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Router interface {
	http.Handler
	GetChiRouter() chi.Router
	Route(pattern string, fn func(r Router)) Router
	SetAuthMiddleware(middlewares ...func(http.Handler) http.Handler)
	SetOptAuthMiddleware(middlewares ...func(http.Handler) http.Handler)
	Use(middlewares ...func(http.Handler) http.Handler)
	With(middlewares ...func(http.Handler) http.Handler) Router
	Auth() Router
	OptAuth() Router
	Connect(pattern string, re RouterElement)
	Delete(pattern string, re RouterElement)
	Get(pattern string, re RouterElement)
	Head(pattern string, re RouterElement)
	Options(pattern string, re RouterElement)
	Patch(pattern string, re RouterElement)
	Post(pattern string, re RouterElement)
	Put(pattern string, re RouterElement)
	Trace(pattern string, re RouterElement)
	// router のエンドポイントと input, output の型定義を出力する
	GetRouterDefinition() ([]*RouterDefinition, map[string]*TypeStructure)
}

type RouterDefinition struct {
	InputTypeStructure  *TypeStructure `json:"input_type_structure"`
	OutputTypeStructure *TypeStructure `json:"output_type_structure"`
	FullPathName        string         `json:"full_path_name"`
	CurrentPathName     string         `json:"current_path_name"`
	Method              string         `json:"method"`
	WithAuth            bool           `json:"with_auth"`
}
