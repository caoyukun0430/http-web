package engine

import (
	"fmt"
	"net/http"
)

// Improvment day2 HandlerFunc takes Context as argument, and engine is still an implementation
// of Handler interface. But we abstract router into its file so that engine is independent
// from router map

// HandlerFunc defines the request handler used
type HandlerFunc func(*Context)

// Handler is an interface that requires any type that implements it to have a ServeHTTP method.
// This method should have the signature ServeHTTP(ResponseWriter, *Request).
// It's a way to create a contract that any HTTP handler must satisfy
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Engine implement the interface of ServeHTTP, difference is HandlerFunc takes Context as argument
type Engine struct {
	router *router
}

// New is the constructor of Engine, init the router map
func New() *Engine {
	return &Engine{router: newRouter()}
}

// implement the Handler interface as Engine pointer as we need to modify Engine map
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handleRoute(c)
}

// add routing rules to router table
func (engine *Engine) AddRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.AddRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) Get(pattern string, handler HandlerFunc) {
	engine.AddRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) Post(pattern string, handler HandlerFunc) {
	engine.AddRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	fmt.Printf("HTTP server starting at %s ...\n", addr)
	return http.ListenAndServe(addr, engine)
}
