package main

import (
	"log"
	"net/http"
	"time"

	"engine"
)

// encapsulate w, r into context like Gin to simplify API interface calls, to avoid repetitive
// construction of http header type, status ... and body. Morever, extend supports of HTML,JSON

func onlyForV2() engine.HandlerFunc {
	return func(c *engine.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	server := engine.New()
	server.AppendMid(engine.Logger()) // global midlleware
	server.Get("/", func(c *engine.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := server.Group("/v2")
	v2.AppendMid(onlyForV2()) // v2 group middleware
	{
		v2.Get("/hello/:name", func(c *engine.Context) {
			// expect /hello/geektutu
			c.Plain(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	server.Run(":8080")
}
