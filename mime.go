// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

package glaze

const (
	MIME_JSON                = "application/json"
	MIME_HTML                = "text/html"
	MIME_PLAIN               = "text/plain"
	MIME_POST_FORM           = "application/x-www-form-urlencoded"
	MIME_MULTIPART_POST_FORM = "multipart/form-data"
)

var (
	jsonContentType      = []string{"application/json; charset=utf-8"}
	textPlainContentType = "text/plain; charset=utf-8"
)
