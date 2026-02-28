package router

import (
	"context"
	"net/http"
	"strings"
)

type paramsKey struct{}

type Router struct {
	routes []route
}

type route struct {
	method   string
	pattern  string
	segments []string
	handler  http.Handler
}

func New() *Router { return &Router{} }

func (r *Router) Handle(method, pattern string, h http.Handler) {
	r.routes = append(r.routes, route{
		method:   method,
		pattern:  pattern,
		segments: splitPath(pattern),
		handler:  h,
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathSeg := splitPath(req.URL.Path)

	for _, rt := range r.routes {
		if rt.method != req.Method {
			continue
		}
		if len(rt.segments) != len(pathSeg) {
			continue
		}

		params := map[string]string{}
		matched := true

		for i := range rt.segments {
			p := rt.segments[i]
			s := pathSeg[i]

			if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
				key := strings.TrimSuffix(strings.TrimPrefix(p, "{"), "}")
				params[key] = s
				continue
			}
			if p != s {
				matched = false
				break
			}
		}

		if matched {
			ctx := context.WithValue(req.Context(), paramsKey{}, params)
			rt.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}

	http.NotFound(w, req)
}

func Param(r *http.Request, name string) string {
	m, _ := r.Context().Value(paramsKey{}).(map[string]string)
	return m[name]
}

func splitPath(p string) []string {
	p = strings.TrimSpace(p)
	if p == "" || p == "/" {
		return []string{}
	}
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}
