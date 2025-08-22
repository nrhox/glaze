// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

const defaultMultipartMemory = 40 << 20 // default size 40 MB

// Engine is the main object for the web framework.
// It holds routes, configs, trees, and HTTP server features.
type Engine struct {
	Route
	routeList   []RouteInfo // all routes information
	releaseMode bool        // flag for release mode

	writer          io.Writer        // where log is written
	MultipartMemory int64            // memory limit for multipart form
	trees           map[string]*node // route trees (per method)
}

// make sure Engine implement IRouter
var _ Router = (*Engine)(nil)

// ConfigsFunc is a function type that change engine configuration.
type ConfigsFunc func(*Engine)

// New create new engine with default configs.
// Can pass options function to change config.
func New(cfg ...ConfigsFunc) *Engine {
	engine := &Engine{
		MultipartMemory: defaultMultipartMemory,
		trees:           make(map[string]*node),
		writer:          os.Stdout,
	}

	// self reference to engine
	engine.engine = engine
	return engine.Config(cfg...)
}

// Config apply all config functions to engine.
func (e *Engine) Config(cfgs ...ConfigsFunc) *Engine {
	for _, opt := range cfgs {
		opt(e) // run option function
	}
	return e
}

// RoutesInfo return all routes info sorted by path length.
// Useful for debug or listing routes.
func (e *Engine) RoutesInfo() []RouteInfo {
	result := make([]RouteInfo, len(e.routeList))
	copy(result, e.routeList)

	// sort: first by length, then alphabet
	sort.Slice(result, func(i, j int) bool {
		if len(result[i].Path) == len(result[j].Path) {
			return result[i].Path < result[j].Path
		}
		return len(result[i].Path) > len(result[j].Path)
	})
	return result
}

// ServeHTTP implement http.Handler.
// It find route, create context, and run handlers.
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handlers, params := e.findRoute(req.Method, req.URL.Path)
	if handlers == nil {
		// if route not found, return 404
		http.NotFound(w, req)
		return
	}

	// create context for this request
	c := &Context{
		Writer:   w,
		Request:  req,
		Params:   params,
		handlers: handlers,
		index:    -1,
		querys:   req.URL.Query(),
		engine:   e.engine,
	}
	// start handler chain
	c.Next()
}

// RunAndListen starts an HTTP server at the given address.
// This function is simple: it does not support graceful shutdown.
//
// Example:
//
//	e := glaze.New()
//	e.RunAndListen(":8080")
func (e *Engine) RunAndListen(addr string) error {
	if !e.releaseMode {
		// show all routes in console
		for _, r := range e.RoutesInfo() {
			fmt.Fprintf(e.writer, "%-6s %s\n", r.Method, r.Path)
		}
	}
	if !e.releaseMode {
		fmt.Fprintf(e.writer, "listen on %s\n", addr)
	}
	return http.ListenAndServe(addr, e)
}

// ListenAndGraceful starts an HTTP server at the given address,
// but it also listen for system signals (SIGINT, SIGTERM).
// When signal received, it shutdown the server gracefully with timeout.
//
// Example:
//
//	e := glaze.New()
//	e.ListenAndGraceful(":8080")
func (e *Engine) ListenAndGraceful(addr string) error {
	if !e.releaseMode {
		// show all routes in console
		for _, r := range e.RoutesInfo() {
			fmt.Fprintf(e.writer, "%-6s %s\n", r.Method, r.Path)
		}
	}

	// create http server
	srv := &http.Server{
		Addr:    addr,
		Handler: e.engine,
	}

	if !e.releaseMode {
		fmt.Fprintf(e.writer, "listen on %s\n", addr)
	}

	// run server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(e.writer, "listen: %s\n", err)
		}
	}()

	// wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Fprint(e.writer, "Shutdown Server")

	// graceful shutdown with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	fmt.Fprint(e.writer, "Server exiting")
	return nil
}
