package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H is a shortcut for map[string]interface{} like gin
//
//	obj = map[string]interface{}{
//	    "name": "geektutu",
//	    "password": "1234",
//	}
type H map[string]interface{}

// day5 update: middleware, it's similar/kind of HandlerFunc that it context targeted and can be
// defined by user to record logs, calculate latency, ... Middleware is
// appended after the request is received and context is initialized.
// By defining Next() func inside cotext, we construct the stack-struct FILO middleware-chain,
// and each handler has freedom to choice to call Next() to split the work pre request and post
// request.
// func middlewareA(c *Context) {
//     // first part to handle pre request
//     part1
//     c.Next()
//     // sec part to handle pre request
//     part2
// }

// update: introduce Params to map pattern-path for wildcard pattern
// we use context to encapsulate *http.Request and http.ResponseWriter
// so that we can hide the implementation details about header, body
// and different types, html, json, plain text ...
type Context struct {
	Req    *http.Request
	Writer http.ResponseWriter
	// request info
	Method string
	Path   string
	Params map[string]string
	// resp
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	engine   *Engine // engine pointer used in HTML
}

// constructor func
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Req:    r,
		Writer: w,
		Method: r.Method,
		Path:   r.URL.Path,
		index:  -1,
	}
}

// Next() maintains middleware stack
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail works as curcuit breaker, when called, all handlers after are skipped
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// basic method for wildcard Params
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// basic methods FormValue and Query
// FormValue returns the first value for the named component of the query.
// POST and PUT body parameters take precedence over URL query string values.
// FormValue calls ParseMultipartForm and ParseForm if necessary and ignores
// any errors returned by these functions.
// If key is not present, FormValue returns the empty string.
// To access multiple values of the same key, call ParseForm and
// then inspect Request.Form directly.

// methods has pointer receiver so that we can modify Context instance
func (c *Context) FormValue(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	// WriteHeader sends an HTTP response header with the provided
	// status code.
	c.Writer.WriteHeader(code)
}

// src\net\http\header.go
// Set sets the header entries associated with key to the
// single element value. It replaces any existing values
// associated with key.
func (c *Context) SetHeader(key string, val string) {
	c.Writer.Header().Set(key, val)
}

// three types plain-text, HTML, JSON

// it means you can pass in a list of arguments of any type, and those arguments
// will be accessible within the function as a slice of interface{}.
// The function fmt.Sprintf is then used to format those values according to the specified format string.
// c.Plain(http.StatusOK, "User %s has %d messages", "Alice", 25)
// Alice" and 25 are passed as arguments to values ...interface{} and will be formatted into the string by fmt.Sprintf
func (c *Context) Plain(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// refer to https://github.com/gin-gonic/gin/blob/master/render/json.go#L179
// recommend to panic if Encode fails
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

// day6 improve HMTL method so that it can render base on template name and data received
func (c *Context) HTML(code int, tmpl string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, tmpl, data); err != nil {
		c.Fail(500, err.Error())
	}
}
