package main

import (
	"net/http"

	"engine"
)

// encapsulate w, r into context like Gin to simplify API interface calls, to avoid repetitive
// construction of http header type, status ... and body. Morever, extend supports of HTML,JSON

func main() {
	server := engine.New()
	server.Get("/index", func(c *engine.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := server.Group("/v1")
	{
		v1.Get("/", func(c *engine.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.Get("/hello", func(c *engine.Context) {
			// expect /hello?name=geektutu
			c.Plain(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := server.Group("/v2")
	{
		v2.Get("/hello/:name", func(c *engine.Context) {
			// expect /hello/geektutu
			c.Plain(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.Post("/login", func(c *engine.Context) {
			c.JSON(http.StatusOK, engine.H{
				"username": c.FormValue("username"),
				"password": c.FormValue("password"),
			})
		})

	}

	server.Run(":8080")
}
