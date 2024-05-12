package main

import (
	"net/http"

	"engine"
)

// encapsulate w, r into context like Gin to simplify API interface calls, to avoid repetitive
// construction of http header type, status ... and body. Morever, extend supports of HTML,JSON

func main() {
	server := engine.New()
	server.Get("/", func(c *engine.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Web</h1>")
	})
	server.Get("/hello", func(c *engine.Context) {
		// expect /hello?name=geektutu
		c.Plain(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	server.Post("/login", func(c *engine.Context) {
		c.JSON(http.StatusOK, engine.H{
			"username": c.FormValue("username"),
			"password": c.FormValue("password"),
		})
	})

	server.Run(":8080")
}
