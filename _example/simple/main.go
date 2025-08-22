// Copyright 2025 Jalu Nugroho
// SPDX-License-Identifier: MIT

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
}
