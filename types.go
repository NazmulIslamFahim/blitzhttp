package blitzhttp

import (
	"net/http"
)

// Middleware wraps a handler
type Middleware func(http.Handler) http.Handler

// paramsKey stores route parameters in context
type paramsKeyType struct{}

var paramsKey = paramsKeyType{}

// GetParams retrieves route parameters
func GetParams(r *http.Request) string {
	if params, ok := r.Context().Value(paramsKey).(string); ok {
		return params
	}
	return ""
}

// method indices for switch-based dispatch
const (
	mGet = iota
	mPost
	mPut
	mDelete
	mPatch
	mAny
)

// routeHandler represents a route handler
type routeHandler struct {
	handler http.Handler // Pre-composed with middlewares
}
