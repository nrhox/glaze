// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

import (
	"net/http"
	"testing"
)

func BenchmarkOneRoute(B *testing.B) {
	router := New()
	router.Get("/ping", func(c *Context) {})
	runRequest(B, router, http.MethodGet, "/ping")
}

func Benchmark5Params(B *testing.B) {
	router := New()
	router.Use(func(c *Context) {})
	router.Get("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {})
	runRequest(B, router, http.MethodGet, "/param/path/to/parameter/john/12345")
}

func Benchmark10Params(B *testing.B) {
	router := New()
	router.Use(func(c *Context) {})
	router.Get("/param/:param1/:params2/:param3/:param4/:param5/:param6/:param7/:param8/:param9/:param10", func(c *Context) {})
	runRequest(B, router, http.MethodGet, "/param/horeg/with/test/performance/and/wish/performance/lessthen/700ns")
}

type mockRequest struct {
	headers http.Header
}

func newMockRequest() *mockRequest {
	return &mockRequest{
		http.Header{},
	}
}

func (m *mockRequest) Header() (h http.Header) {
	return m.headers
}

func (m *mockRequest) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockRequest) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockRequest) WriteHeader(int) {}

func runRequest(B *testing.B, r *Engine, method, path string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := newMockRequest()
	B.ReportAllocs()

	for B.Loop() {
		r.ServeHTTP(w, req)
	}
}
