package rapi

import (
	"net/http"

	"github.com/go-chi/chi"
)

func NewRouter() Router {
	r := &router{
		chiRouter:          chi.NewRouter(),
		children:           []*router{},
		authMiddlewares:    chi.Middlewares{},
		optAuthMiddlewares: chi.Middlewares{},
	}
	r.root = r
	return r
}

type router struct {
	method             string
	path               string
	root               *router
	parent             *router
	withAuth           bool
	chiRouter          chi.Router
	element            RouterElement
	children           []*router
	authMiddlewares    chi.Middlewares
	optAuthMiddlewares chi.Middlewares
}

func (r *router) sub() *router {
	subRouter := &router{
		children: []*router{},
		root:     r.root,
		parent:   r,
	}
	r.children = append(r.children, subRouter)
	return subRouter
}

func (r *router) handle(method string, pattern string, re RouterElement) {
	subRouter := r.sub()
	subRouter.chiRouter = r.chiRouter
	subRouter.element = re
	subRouter.path = pattern
	subRouter.method = method

	switch method {
	case http.MethodConnect:
		subRouter.chiRouter.Connect(pattern, re.GetHandleFunc())
	case http.MethodDelete:
		subRouter.chiRouter.Delete(pattern, re.GetHandleFunc())
	case http.MethodGet:
		subRouter.chiRouter.Get(pattern, re.GetHandleFunc())
	case http.MethodHead:
		subRouter.chiRouter.Head(pattern, re.GetHandleFunc())
	case http.MethodOptions:
		subRouter.chiRouter.Options(pattern, re.GetHandleFunc())
	case http.MethodPatch:
		subRouter.chiRouter.Patch(pattern, re.GetHandleFunc())
	case http.MethodPost:
		subRouter.chiRouter.Post(pattern, re.GetHandleFunc())
	case http.MethodPut:
		subRouter.chiRouter.Put(pattern, re.GetHandleFunc())
	case http.MethodTrace:
		subRouter.chiRouter.Trace(pattern, re.GetHandleFunc())
	}
}

func (r *router) GetChiRouter() chi.Router {
	return r.chiRouter
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chiRouter.ServeHTTP(w, req)
}

func (r *router) Route(pattern string, fn func(r Router)) Router {
	subRouter := r.sub()
	subRouter.path = pattern

	if fn != nil {
		r.chiRouter.Route(pattern, func(chiRouter chi.Router) {
			subRouter.chiRouter = chiRouter
			fn(subRouter)
		})
	} else {
		subRouter.chiRouter = r.chiRouter.Route(pattern, nil)
	}

	return subRouter
}

func (r *router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.chiRouter.Use(middlewares...)
}

func (r *router) SetAuthMiddleware(middlewares ...func(http.Handler) http.Handler) {
	r.root.authMiddlewares = append(r.authMiddlewares, middlewares...)
}

func (r *router) SetOptAuthMiddleware(middlewares ...func(http.Handler) http.Handler) {
	r.root.optAuthMiddlewares = append(r.optAuthMiddlewares, middlewares...)
}

func (r *router) With(middlewares ...func(http.Handler) http.Handler) Router {
	return r.with(middlewares...)
}

func (r *router) with(middlewares ...func(http.Handler) http.Handler) *router {
	subRouter := r.sub()
	subRouter.chiRouter = r.chiRouter.With(middlewares...)
	return subRouter
}

func (r *router) Auth() Router {
	subRouter := r.with(r.root.authMiddlewares...)
	subRouter.withAuth = true
	return subRouter
}

func (r *router) OptAuth() Router {
	subRouter := r.with(r.root.optAuthMiddlewares...)
	subRouter.withAuth = true
	return subRouter
}

func (r *router) Connect(pattern string, re RouterElement) {
	r.handle(http.MethodConnect, pattern, re)
}

func (r *router) Delete(pattern string, re RouterElement) {
	r.handle(http.MethodDelete, pattern, re)
}

func (r *router) Get(pattern string, re RouterElement) {
	r.handle(http.MethodGet, pattern, re)
}

func (r *router) Head(pattern string, re RouterElement) {
	r.handle(http.MethodHead, pattern, re)
}

func (r *router) Options(pattern string, re RouterElement) {
	r.handle(http.MethodOptions, pattern, re)
}

func (r *router) Patch(pattern string, re RouterElement) {
	r.handle(http.MethodPatch, pattern, re)
}

func (r *router) Post(pattern string, re RouterElement) {
	r.handle(http.MethodPost, pattern, re)
}

func (r *router) Put(pattern string, re RouterElement) {
	r.handle(http.MethodPut, pattern, re)
}

func (r *router) Trace(pattern string, re RouterElement) {
	r.handle(http.MethodTrace, pattern, re)
}

// router のエンドポイントと input, output の型定義を出力する
func (r *router) GetRouterDefinition() ([]*RouterDefinition, map[string]*TypeStructure) {
	ts := NewTypeScanner()
	ts.DisableStructField()
	ts.AddStructTagName("json", "form")

	routerDefinitions := []*RouterDefinition{}

	// 再帰で全てのRouter定義をappend
	var appendRouterDefinition func(r *router, parentPath string)
	appendRouterDefinition = func(r *router, parentPath string) {
		if r.element != nil {
			routerDefinition := &RouterDefinition{
				FullPathName:        parentPath + r.path,
				CurrentPathName:     r.path,
				Method:              r.method,
				WithAuth:            r.withAuth,
				InputTypeStructure:  ts.Scan(r.element.GetEmptyInput()),
				OutputTypeStructure: ts.Scan(r.element.GetEmptyOutput()),
			}
			routerDefinitions = append(routerDefinitions, routerDefinition)
		}

		for _, child := range r.children {
			appendRouterDefinition(child, parentPath+r.path)
		}
	}

	appendRouterDefinition(r.root, "")
	return routerDefinitions, ts.Export()
}
