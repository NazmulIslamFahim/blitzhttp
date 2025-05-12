package blitzhttp

import (
	"context"
	"net/http"
	"strings"
)

// Router is the main router
type Router struct {
	static     map[string]*route                        // Exact path routes
	params     map[string]*route                        // Parameterized routes (e.g., /users/:id)
	wildcards  map[string]*route                        // Wildcard routes (e.g., /files/*)
	catchAll   *route                                   // Catch-all route (e.g., *)
	globalMWs  []Middleware                             // Global middlewares
	staticCode func(http.ResponseWriter, *http.Request) // Compiled static routes
}

// route represents a path's handlers
type route struct {
	handlers [mAny + 1]*routeHandler // Method-specific handlers
	isParam  bool                    // Has :param
	isWild   bool                    // Has * wildcard
}

// New creates a Router
func New() *Router {
	r := &Router{
		static:    make(map[string]*route),
		params:    make(map[string]*route),
		wildcards: make(map[string]*route),
	}
	r.compileStatic()
	return r
}

// ServeHTTP handles requests
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		w.Header().Set("Allow", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	if r.staticCode != nil {
		r.staticCode(w, req)
		return
	}

	// Exact match
	if rt, ok := r.static[path]; ok {
		if h := rt.getHandler(req.Method); h != nil {
			h.handler.ServeHTTP(w, req)
			return
		}
	}

	// // Parameterized match
	// for p, rt := range r.params {
	// 	if params := matchParam(path, p); params != "" {
	// 		if h := rt.getHandler(req.Method); h != nil {
	// 			ctx := context.WithValue(req.Context(), paramsKey, params)
	// 			h.handler.ServeHTTP(w, req.WithContext(ctx))
	// 			return
	// 		}
	// 	}
	// }

	// // Wildcard match
	// for p, rt := range r.wildcards {
	// 	if strings.HasPrefix(path, p[:len(p)-1]) {
	// 		if h := rt.getHandler(req.Method); h != nil {
	// 			params := path[len(p):]
	// 			ctx := context.WithValue(req.Context(), paramsKey, params)
	// 			h.handler.ServeHTTP(w, req.WithContext(ctx))
	// 			return
	// 		}
	// 	}
	// }

	// Catch-all match
	if r.catchAll != nil {
		if h := r.catchAll.getHandler(req.Method); h != nil {
			ctx := context.WithValue(req.Context(), paramsKey, path)
			h.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}

	http.NotFound(w, req)
}

// matchParam checks if path matches a parameterized pattern
func matchParam(path, pattern string) string {
	if len(path) < len(pattern) {
		return ""
	}
	if pattern == path {
		return ""
	}
	pParts := strings.Split(pattern, "/")
	rParts := strings.Split(path, "/")
	if len(pParts) != len(rParts) {
		return ""
	}
	for i, p := range pParts {
		if p == rParts[i] || p[0] == ':' {
			continue
		}
		return ""
	}
	return strings.Join(rParts, "/")
}

// Use adds global middlewares
func (r *Router) Use(mws ...Middleware) {
	r.globalMWs = append(r.globalMWs, mws...)
	r.recomposeHandlers()
}

// recomposeHandlers updates all handlers with global middlewares
func (r *Router) recomposeHandlers() {
	for _, routes := range []map[string]*route{r.static, r.params, r.wildcards} {
		for _, rt := range routes {
			for i, h := range rt.handlers {
				if h != nil {
					rt.handlers[i].handler = composeHandler(h.handler, r.globalMWs...)
				}
			}
		}
	}
	if r.catchAll != nil {
		for i, h := range r.catchAll.handlers {
			if h != nil {
				r.catchAll.handlers[i].handler = composeHandler(h.handler, r.globalMWs...)
			}
		}
	}
	r.compileStatic()
}
