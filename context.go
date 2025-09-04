// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

import (
	"encoding/json"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

// Context is like the request context in web framework.
// It hold request, response, params, query, handlers, and custom values.
type Context struct {
	Writer  http.ResponseWriter // write response back
	Request *http.Request       // http request
	Params  map[string]string   // path parameters like /user/:id
	querys  url.Values          // query parameters

	handlers []HandlerFunc // list of handler functions (middlewares)
	index    int           // current handler index
	engine   *Engine       // pointer to engine

	Keys map[any]any  // custom key-value storage
	mu   sync.RWMutex // lock for safe access

	stopped bool // stop flag to abort next handlers
}

// Next call the next handler in the list.
// It move index and run handler one by one.
func (c *Context) Next() {
	if c.stopped {
		return // if stopped, no continue
	}
	// move to next handler
	c.index++

	// check if still inside handler list
	if c.index < len(c.handlers) {
		// call handler function
		c.handlers[c.index](c)
		// continue to next
		c.Next()
	}
}

// Abort stop the handler chain.
// After this, no more handler will be run.
func (c *Context) Abort() {
	// just set stop flag
	c.stopped = true
}

// Param return value from path parameter by key.
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// Query return value from query parameter in URL.
func (c *Context) Query(key string) string {
	return c.querys.Get(key)
}

// Set put a custom value inside context.
func (c *Context) Set(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Keys == nil {
		c.Keys = make(map[any]any)
	}
	c.Keys[key] = value
}

// Get return a custom value from context.
func (c *Context) Get(key any) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.Keys[key]
	return
}

// String write plain text response with status code.
func (c *Context) String(code int, msg string) {
	c.Writer.WriteHeader(code)
	io.WriteString(c.Writer, msg)
}

// writeContentType set Content-Type header if not exist.
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// JSON send JSON response with escape HTML off.
func (c *Context) JSON(code int, data any) {
	writeContentType(c.Writer, jsonContentType)
	c.Writer.WriteHeader(code)

	encoder := json.NewEncoder(c.Writer)
	encoder.SetEscapeHTML(false) // do not escape HTML
	encoder.Encode(data)
}

// PureJSON send JSON response with escape HTML on.
func (c *Context) PureJSON(code int, data any) {
	writeContentType(c.Writer, jsonContentType)
	c.Writer.WriteHeader(code)

	encoder := json.NewEncoder(c.Writer)
	encoder.SetEscapeHTML(true) // escape HTML
	encoder.Encode(data)
}

// BindJSON read JSON request body and decode into struct.
func (c *Context) BindJSON(dst any) error {
	if c.Request.Header.Get("Content-Type") != MIME_JSON {
		return http.ErrNotSupported
	}

	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(dst); err != nil {
		return err
	}
	return nil
}

// FormFile return uploaded file header by field name.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(c.engine.MultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// MultipartForm return multipart form data from request.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(c.engine.MultipartMemory)
	return c.Request.MultipartForm, err
}

// SaveFile save uploaded file to destination path.
func (c *Context) SaveFile(file *multipart.FileHeader, dst string, perm ...fs.FileMode) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var mode os.FileMode = 0o750
	if len(perm) > 0 {
		mode = perm[0]
	}
	dir := filepath.Dir(dst)
	if err = os.MkdirAll(dir, mode); err != nil {
		return err
	}
	if err = os.Chmod(dir, mode); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// SetCookie add a cookie into response.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure bool, httpOnly bool, sameSite http.SameSite) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	})
}

// GetCookie return cookie value from request by name.
func (c *Context) GetCookie(key string) (string, error) {
	cookie, err := c.Request.Cookie(key)
	if err != nil {
		return "", err
	}
	v, _ := url.QueryUnescape(cookie.Value)
	return v, nil
}

// GetHeader return header value by key from request
func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// set header into the runtime request or response
func (c *Context) SetHeader(key, value string) {
	c.Request.Header.Set(key, value)
}
