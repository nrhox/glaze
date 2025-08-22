// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	msg := "Hello World"

	r := New()
	r.Get("/ping", func(c *Context) {
		c.String(200, msg)
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "they should be equal")
	assert.Equal(t, msg, string(body), "they should be equal")
}

func TestHttpWithParam(t *testing.T) {
	msgParam := "Hello_World"

	r := New()

	r.Get("/p/:param", func(c *Context) {
		param := c.Param("param")
		c.String(200, param)
	})

	req := httptest.NewRequest("GET", "/p/"+msgParam, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "they should be equal")
	assert.Equal(t, msgParam, string(body), "they should be equal")
}

func TestHttpWithQuery(t *testing.T) {
	query := "car blue"

	r := New()

	r.Get("/p", func(c *Context) {
		keyword := c.Query("q")

		c.String(200, keyword)
	})

	req := httptest.NewRequest("GET", "/p?q="+url.QueryEscape(query), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "they should be equal")
	assert.Equal(t, query, string(body), "they should be equal")
}

func auth() HandlerFunc {
	return func(c *Context) {
		if p := c.Param("id"); p != "Hello" {
			c.String(http.StatusNotFound, "param should be Hello")
			c.Abort()
			return
		}

	}
}

func TestHttpWithMiddleware(t *testing.T) {
	r := New()

	r.Get("/p/:id", auth(), func(c *Context) {
		keyword := c.Query("q")

		c.String(200, keyword)
	})

	req := httptest.NewRequest("GET", "/p/ds", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "they should be equal")
	assert.Equal(t, "param should be Hello", string(body), "they should be equal")
}

func TestHttpWithMiddlewareSuccess(t *testing.T) {
	r := New()

	r.Get("/p/:id", auth(), func(c *Context) {
		keyword := c.Query("q")

		c.String(200, keyword)
	})

	req := httptest.NewRequest("GET", "/p/Hello?q=world", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode, "they should be equal")
	assert.Equal(t, "world", string(body), "they should be equal")
}
