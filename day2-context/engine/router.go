package engine

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// add routing rules to router table
func (r *router) AddRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) handleRoute(c *Context) {
	key := c.Method + "-" + c.Path
	handler, ok := r.handlers[key]
	// If the key exists
	if ok {
		handler(c)
	} else {
		c.Plain(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
