// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

import (
	"net/http"
	"path"
	"regexp"
)

// HandlerFunc defines a request handler used by the framework.
// Every handler receives a Context, which holds request and response data.
type HandlerFunc func(*Context)

// HandlersChain represents a list of HandlerFunc executed in order.
// This is typically used for middleware + final handler.
type HandlersChain []HandlerFunc

// M is a shortcut for building JSON objects.
// Example: c.JSON(200, glaze.M{"msg": "ok"})
type M map[string]any

// RouteInfo describes a single registered route,
// including the HTTP method and the route path.
type RouteInfo struct {
	Method string
	Path   string
}

// Router is the main interface for grouping and
// registering routes. It embeds Routes for HTTP
// method helpers and adds Group for sub-routes.
type Router interface {
	Routes
	Group(string, ...HandlerFunc) *Route
}

// Routes defines the basic routing methods available
// for each HTTP method, and also allows attaching
// middleware with Use.
type Routes interface {
	Use(...HandlerFunc) Routes

	Get(string, ...HandlerFunc) Routes
	Post(string, ...HandlerFunc) Routes
	Delete(string, ...HandlerFunc) Routes
	Patch(string, ...HandlerFunc) Routes
	Put(string, ...HandlerFunc) Routes
	Options(string, ...HandlerFunc) Routes
	Head(string, ...HandlerFunc) Routes
}

// Route represents a registered route or a route group.
// It stores the HTTP method, path, and handlers chain.
// Nested groups keep track of their parent engine.
type Route struct {
	Method  string
	Path    string
	Handler HandlersChain
	root    bool
	engine  *Engine
}

// ensure Route implements IRouter
var _ Router = (*Route)(nil)

var (
	regexMethodLetter = regexp.MustCompile("^[A-Z]+$")
)

// Use appends middleware handlers to the route or group.
func (r *Route) Use(middleware ...HandlerFunc) Routes {
	r.Handler = append(r.Handler, middleware...)
	return r.engineInfo()
}

// Group creates a new route group with a common path prefix
// and optional middleware handlers.
func (r *Route) Group(path string, handlers ...HandlerFunc) *Route {
	return &Route{
		engine:  r.engine,
		Path:    r.jointAbsolutePath(path),
		Handler: r.joinHandler(handlers),
	}
}

// handle registers a new route with the given HTTP method,
// relative path, and handlers.
func (r *Route) handle(method, relativePath string, handlers ...HandlerFunc) Routes {
	if matched := regexMethodLetter.MatchString(method); !matched {
		panic("invalid method '" + method + "'")
	}
	absolutePath := r.jointAbsolutePath(relativePath)
	handlers = r.joinHandler(handlers)
	r.engine.addRoute(method, absolutePath, handlers...)
	return r.engineInfo()
}

func (r *Route) Get(path string, handler ...HandlerFunc) Routes {
	return r.handle(http.MethodGet, path, handler...)
}
func (r *Route) Post(path string, handler ...HandlerFunc) Routes {
	return r.handle(http.MethodPost, path, handler...)
}
func (r *Route) Put(path string, handler ...HandlerFunc) Routes {
	return r.handle(http.MethodPut, path, handler...)
}
func (r *Route) Delete(path string, handler ...HandlerFunc) Routes {
	return r.handle(http.MethodDelete, path, handler...)
}
func (r *Route) Patch(path string, handler ...HandlerFunc) Routes {
	return r.handle(http.MethodPatch, path, handler...)
}
func (r *Route) Options(path string, handler ...HandlerFunc) Routes {
	return r.handle(http.MethodOptions, path, handler...)
}
func (r *Route) Head(path string, handler ...HandlerFunc) Routes {
	return r.handle(http.MethodHead, path, handler...)
}

// joinHandler merges current handlers with new handlers,
// preserving order (middleware first, then final handler).
func (r *Route) joinHandler(handlers HandlersChain) HandlersChain {
	finalSize := len(r.Handler) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, r.Handler)
	copy(mergedHandlers[len(r.Handler):], handlers)
	return mergedHandlers
}

// jointAbsolutePath combines the parent route path with a child path.
func (r *Route) jointAbsolutePath(relativePath string) string {
	return joinPath(r.Path, relativePath)
}

// lastChar returns the last character of a string, or panics if empty.
func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

// joinPath combines absolute and relative paths
// and ensures trailing slash is preserved if present.
func joinPath(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

// engineInfo returns the Routes reference,
// pointing back to the engine if this is the root group.
func (group *Route) engineInfo() Routes {
	if group.root {
		return group.engine
	}
	return group
}
