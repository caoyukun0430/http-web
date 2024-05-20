package engine

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Improvment day2 HandlerFunc takes Context as argument, and engine is still an implementation
// of Handler interface. But we abstract router into its file so that engine is independent
// from router map

// Improvement day4 Add group management feature for routing based on prefix. So that we can controll
// nested routing pre prefix group and later will be convenient to add middlewares to routing groups
// we want. To achieve, all groups share the same Engine by including the Engine pointer into the
// RouterGroup struct

// HandlerFunc defines the request handler used
type HandlerFunc func(*Context)

// Handler is an interface that requires any type that implements it to have a ServeHTTP method.
// This method should have the signature ServeHTTP(ResponseWriter, *Request).
// It's a way to create a contract that any HTTP handler must satisfy
type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// day4 router group
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine // all groups share a Engine instance
}

// Engine implement the interface of ServeHTTP, difference is HandlerFunc takes Context as argument
type Engine struct {
	*RouterGroup //embedded type
	router       *router
	groups       []*RouterGroup // store all groups into engine
}

// New is the constructor of Engine, init the router map
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	// fmt.Printf("group size %d\n", len(engine.groups)) // size = 1
	return engine
}

// Group is defined to create a new RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	// remember all groups share the same Engine instance
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix, // nested routing
		parent: group,                 // parent is the receiver group for nesting
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// add routing rules to router table
func (group *RouterGroup) addRoute(method string, prefix string, handler HandlerFunc) {
	pattern := group.prefix + prefix
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) Get(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) Post(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// AppendMid append middlewares to certain router group
func (group *RouterGroup) AppendMid(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// day5 update when request reached, all middlewares belong to be URL group
// are added before the request handler
// implement the Handler interface as Engine pointer as we need to modify Engine map
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	engine.router.handleRoute(c)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	fmt.Printf("HTTP server starting at %s ...\n", addr)
	return http.ListenAndServe(addr, engine)
}
