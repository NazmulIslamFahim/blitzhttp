package blitzhttp

import (
	"context"
	"net/http"
	"strings"
)

// getHandler returns the handler for a method
func (rt *route) getHandler(method string) *routeHandler {
	switch method {
	case http.MethodGet:
		return rt.handlers[mGet]
	case http.MethodPost:
		return rt.handlers[mPost]
	case http.MethodPut:
		return rt.handlers[mPut]
	case http.MethodDelete:
		return rt.handlers[mDelete]
	case http.MethodPatch:
		return rt.handlers[mPatch]
	}
	return nil
}

// GET registers a GET route
func (r *Router) GET(path string, handler http.HandlerFunc, mws ...Middleware) {
	r.addRoute(mGet, path, handler, mws...)
}

// POST registers a POST route
func (r *Router) POST(path string, handler http.HandlerFunc, mws ...Middleware) {
	r.addRoute(mPost, path, handler, mws...)
}

// PUT registers a PUT route
func (r *Router) PUT(path string, handler http.HandlerFunc, mws ...Middleware) {
	r.addRoute(mPut, path, handler, mws...)
}

// DELETE registers a DELETE route
func (r *Router) DELETE(path string, handler http.HandlerFunc, mws ...Middleware) {
	r.addRoute(mDelete, path, handler, mws...)
}

// PATCH registers a PATCH route
func (r *Router) PATCH(path string, handler http.HandlerFunc, mws ...Middleware) {
	r.addRoute(mPatch, path, handler, mws...)
}

// ANY registers a route for all methods
func (r *Router) ANY(path string, handler http.HandlerFunc, mws ...Middleware) {
	if path == "*" {
		r.addCatchAll(handler, mws...)
		return
	}
	for i := mGet; i <= mPatch; i++ {
		r.addRoute(i, path, handler, mws...)
	}
}

func (r *Router) addRoute(method int, path string, handler http.HandlerFunc, mws ...Middleware) {
	path = strings.Trim(path, "/")
	isWild := strings.HasSuffix(path, "*")
	isParam := strings.Contains(path, ":")
	if isWild {
		path = strings.TrimSuffix(path, "*")
	}

	var routes map[string]*route
	switch {
	case isWild:
		routes = r.wildcards
	case isParam:
		routes = r.params
	default:
		routes = r.static
	}

	if _, exists := routes[path]; !exists {
		routes[path] = &route{}
	}

	rt := routes[path]
	rt.isParam = isParam
	rt.isWild = isWild
	mws = append(r.globalMWs, mws...)
	rt.handlers[method] = &routeHandler{handler: composeHandler(handler, mws...)}
	r.compileStatic()
}

// addCatchAll registers a catch-all route
func (r *Router) addCatchAll(handler http.HandlerFunc, mws ...Middleware) {
	if r.catchAll == nil {
		r.catchAll = &route{}
	}
	mws = append(r.globalMWs, mws...)
	for i := mGet; i <= mPatch; i++ {
		r.catchAll.handlers[i] = &routeHandler{handler: composeHandler(handler, mws...)}
	}
}

// compileStatic generates a static dispatch function
func (r *Router) compileStatic() {
	if len(r.static) == 0 {
		r.staticCode = nil
		return
	}
	r.staticCode = func(w http.ResponseWriter, req *http.Request) {
		path := strings.Trim(req.URL.Path, "/")
		if rt, ok := r.static[path]; ok {
			if h := rt.getHandler(req.Method); h != nil {
				h.handler.ServeHTTP(w, req)
				return
			}
		}
		// Check catch-all for unmatched paths
		if r.catchAll != nil {
			if h := r.catchAll.getHandler(req.Method); h != nil {
				ctx := context.WithValue(req.Context(), paramsKey, path)
				h.handler.ServeHTTP(w, req.WithContext(ctx))
				return
			}
		}
		http.NotFound(w, req)
	}
}

// composeHandler combines handler with middlewares
func composeHandler(handler http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}
