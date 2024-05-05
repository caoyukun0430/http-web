package engine

import (
	"fmt"
	"net/http"
)

// Improvment 2: In order to diff GET/POST, we extend our routing table by string-handler map
// map[string]HandlerFunc where key is METHOD-PATH, and val is HandlerFunc
// key idea is Engine still handles all routing and inside ServeHTTP(), it disributes works
// to individual HandlerFunc in the next level, which is ordinary func to implements Handler
// interface simply

// Engine implement the interface of ServeHTTP
type Engine struct {
	router map[string]http.HandlerFunc
}

// New is the constructor of Engine, init the router map
func New() *Engine {
	return &Engine{router: make(map[string]http.HandlerFunc)}
}

// implement the Handler interface as Engine pointer as we need to modify Engine map
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	handler, ok := engine.router[key]
	// If the key exists
	if ok {
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

// add routing rules to router table
func (engine *Engine) AddRoute(method string, pattern string, handler http.HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET defines the method to add GET request
func (engine *Engine) Get(pattern string, handler http.HandlerFunc) {
	engine.AddRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) Post(pattern string, handler http.HandlerFunc) {
	engine.AddRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	fmt.Printf("HTTP server starting at %s ...\n", addr)
	return http.ListenAndServe(addr, engine)
}
