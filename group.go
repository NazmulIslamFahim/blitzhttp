package blitzhttp

import (
	"net/http"
	"strings"
)

// Group represents a route group
type Group struct {
	router *Router
	prefix string
	mws    []Middleware
}

// Group creates a new route group
func (r *Router) Group(prefix string, mws ...Middleware) *Group {
	return &Group{
		router: r,
		prefix: strings.Trim(prefix, "/"),
		mws:    mws,
	}
}

// GET registers a GET route
func (g *Group) GET(path string, handler http.HandlerFunc, mws ...Middleware) {
	g.addRoute(mGet, path, handler, mws...)
}

// POST registers a POST route
func (g *Group) POST(path string, handler http.HandlerFunc, mws ...Middleware) {
	g.addRoute(mPost, path, handler, mws...)
}

// PUT registers a PUT route
func (g *Group) PUT(path string, handler http.HandlerFunc, mws ...Middleware) {
	g.addRoute(mPut, path, handler, mws...)
}

// DELETE registers a DELETE route
func (g *Group) DELETE(path string, handler http.HandlerFunc, mws ...Middleware) {
	g.addRoute(mDelete, path, handler, mws...)
}

// PATCH registers a PATCH route
func (g *Group) PATCH(path string, handler http.HandlerFunc, mws ...Middleware) {
	g.addRoute(mPatch, path, handler, mws...)
}

// ANY registers a route for all methods
func (g *Group) ANY(path string, handler http.HandlerFunc, mws ...Middleware) {
	if path == "*" {
		g.addCatchAll(handler, mws...)
		return
	}
	for i := mGet; i <= mPatch; i++ {
		g.addRoute(i, path, handler, mws...)
	}
}

func (g *Group) addRoute(method int, path string, handler http.HandlerFunc, mws ...Middleware) {
	mws = append(g.mws, mws...)
	g.router.addRoute(method, g.prefix+"/"+strings.Trim(path, "/"), handler, mws...)
}

func (g *Group) addCatchAll(handler http.HandlerFunc, mws ...Middleware) {
	mws = append(g.mws, mws...)
	g.router.addCatchAll(handler, mws...)
}

// Group creates a nested route group
func (g *Group) Group(prefix string, mws ...Middleware) *Group {
	return &Group{
		router: g.router,
		prefix: g.prefix + "/" + strings.Trim(prefix, "/"),
		mws:    append(g.mws, mws...),
	}
}
