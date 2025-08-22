// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Recovery returns a middleware that recovers from panics
// during request handling. It prevents the server from crashing
// and instead logs the panic stack trace, then responds with
// HTTP 500 (Internal Server Error).
//
// Usage:
//
//	r := glaze.New()
//	r.Use(glaze.Recovery())
//	r.GET("/", func(c *glaze.Context) {
//	    panic("something went wrong")
//	})
//
// If a panic occurs, the middleware will:
// 1. Stop the remaining middleware chain.
// 2. Log the panic message and stack trace to the engine's writer.
// 3. Send a 500 response with "Internal Server Error".
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if r := recover(); r != nil {
				// stop next middleware execution
				c.stopped = true

				// capture stack trace for debugging
				stack := debug.Stack()

				// log panic and stack trace
				fmt.Fprintf(c.engine.writer, "[PANIC] %v\n%s\n", r, stack)

				// send 500 response to client
				h := c.Writer.Header()
				if h.Get("Content-Type") == "" {
					h.Set("Content-Type", textPlainContentType)
				}
				c.Writer.WriteHeader(http.StatusInternalServerError)
				_, _ = c.Writer.Write([]byte("Internal Server Error"))
			}
		}()

		// continue executing next handlers if no panic
		c.Next()
	}
}
