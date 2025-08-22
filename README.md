# Glaze

Glaze is a lightweight and fast HTTP router for Go like gin and chi style.
It provides minimal overhead, clean API, and zero dependency (except stdlib).

## Installation

```bash
go get github.com/nrhox/glaze
```

## Quick Start

```go
package main

import (
    "net/http"
    "github.com/nrhox/glaze"
)

func main() {
    r := glaze.New()

    r.GET("/ping", func(c *glaze.Context) {
        c.Writer.Write([]byte("pong"))
    })

    http.ListenAndServe(":8080", r)
}

```

## Example use

### Simple rest api

```go
package main

import "github.com/nrhox/glaze"

func main() {
	r := glaze.New()

	r.GET("/ping", func(c *glaze.Context) {
		c.JSON(200, glaze.H{
			"message": "ping",
		})
	})

	r.GET("/message/:param", func(c *glaze.Context) {
		param := c.Param("param")

		c.JSON(200, glaze.H{
			"message": param,
		})
	})

	r.GET("/message/:param/:param2", func(c *glaze.Context) {
		param := c.Param("param")
		param2 := c.Param("param2")

		c.JSON(200, glaze.H{
			"message":   param,
			"message_2": param2,
		})
	})

	g := r.Group("/api")

	g.GET("/query", func(c *glaze.Context) {
		keyword := c.Query("q")

		c.JSON(200, glaze.H{
			"message": keyword,
		})
	})

	r.RunAndListen(":4040")

    // or
    // http.ListenAndServe(":3000", r)
}
```

### Middleware function

```go
func AuthMiddleware() HandlerFunc {
	return func(c *Context) {
        authHeader := c.GetHeader("Authorization")
		tokenParts := strings.Split(authHeader, " ")
		if authHeader == "" || len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Abort()
			return
		}

		token := tokenParts[1]
        // proccessing
	}
}
```

## License

MIT License

Copyright (c) 2025 Jalu Nugroho
