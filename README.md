# 7 Days Go Web Framework from Scratch

## Web Framework http-web

http-web is a [gin](https://github.com/gin-gonic/gin)-like framework, which follows the 7-day [gee](https://geektutu.com/post/gee.html) implementation.

- Day 1 - http.Handler Interface Basic [Code](http-web/day1-http-base)
- Day 2 - Design a Flexiable Context [Code](gee-web/day2-context)
- Day 3 - Router with Trie-Tree Algorithm [Code](gee-web/day3-router)
- Day 4 - Group Control [Code](gee-web/day4-group)
- Day 5 - Middleware Mechanism [Code](gee-web/day5-middleware)
- Day 6 - Embeded Template Support [Code](gee-web/day6-template)
- Day 7 - Panic Recover & Make it Robust [Code](gee-web/day7-panic-recover)

## Day 1 - Static Route

What we learnt?

1. Utilize Go net/http package to handle HTTP requests, basically the Handler interface as well as the HandlerFunc() for simplicity.

2. Encapsulate Handler inside the Engine interface, and introduce the key(method+path)-value(handler) map for static routing. User can
define custom HandlerFunc() for specific method and path combination, which will be stored inside our custom Engine handler.

```go
func main() {
	server := engine.New()
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL Path is %q\n", r.URL.Path)
	})
	server.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q]: %q\n", k, v)
		}
	})

	server.Run(":8080")
}
```

## Day 2 - Context


What we learnt?

1. Extract a Context struct to handle all HTTP request and response logic, e.g. header, status code and different request content
types like plain-text, JSON, HTML.

2. Extract URL handling into a Router struct in order to enhance it more conveniently in the furture, e.g. support of dynamic
routing.

```go
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
```

## Day 3 - Router

What we learnt?

1. Extend routing functions to support dynamic routing with /:name or /*. To achieve this, we modify the route table from a
simple map to a trie tree, which parses the URL path into trie nodes.

2. Write unit test for routing rules before running main function, this is important to decouple testing so that we can verify
unit feature before running as a whole.

```go
func main() {
	server := engine.New()
	server.Get("/", func(c *engine.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Web</h1>")
	})
	server.Get("/hello", func(c *engine.Context) {
		// expect /hello?name=tom
		c.Plain(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	server.Get("/hello/:name", func(c *engine.Context) {
		// expect /hello/tom
		c.Plain(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	server.Get("/assets/*filepath", func(c *engine.Context) {
		c.JSON(http.StatusOK, engine.H{"filepath": c.Param("filepath")})
	})

	server.Run(":8080")
}
```

## Day 4 - Group

What we learnt?

1. Add group management feature for routing based on prefix. So that we can controll
nested routing pre prefix group and later will be convenient to add middlewares to routing groups
we want. To achieve, all groups share the same Engine by including the Engine pointer into the
RouterGroup struct.


```go
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
```

## Day 5 - Middleware

What we learnt?
1. middleware, it's similar/kind of HandlerFunc that it context targeted and can be defined
by user to record logs, calculate latency, ... Middleware is
appended after the request is received and context is initialized.
By defining Next() func inside cotext, we construct the stack-struct FILO middleware-chain,
and each handler has freedom to choice to call Next() to split the work pre request and post
request.

```go
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
// onlyForV2 executed before logger due to stack struct
$ go run .
2024/05/20 09:49:40 Route  GET - /
2024/05/20 09:49:40 Route  GET - /v2/hello/:name
HTTP server starting at :8080 ...
2024/05/20 09:50:14 [500] /v2/hello/geektutu in 267µs for group v2
2024/05/20 09:50:14 [500] /v2/hello/geektutu in 389.6µs
exit status 0xc000013a
```

## Day 6 - HTTP template

What we learnt?
1. use http/template library to render on the server side with static files.
Static files are mapped with dynamic routing /*filepath. Two rendering functions
SetFuncMap and LoadHTMLGlob are defined for users to define custom render
functions and load template files.

```go
type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := engine.New()
	r.AppendMid(engine.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// load all .tmpl into engine
	r.LoadHTMLGlob("templates/*")
	// map .static/ to /assets URL pattern
	r.Static("/assets", "./static")

	// by default it renders css.tmpl
	r.Get("/", func(c *engine.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.Get("/students", func(c *engine.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", engine.H{
			"title":  "engine",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.Get("/date", func(c *engine.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", engine.H{
			"title": "engine",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":8080")
```

## Day 7 - Error Recovery and Trace
What we learnt?
1. Go panic recover mechanism and we wrote custom trace function to
print stack trace when program panics for debugging in our middleware.

```go
// 2024/06/07 10:03:46 Route  GET - /
// 2024/06/07 10:03:46 Route  GET - /panic
// HTTP server starting at :8080 ...
// 2024/06/07 10:03:59 runtime error: index out of range [100] with length 1
// Traceback:      C:/Program Files/Go/src/runtime/panic.go:884
//         C:/Program Files/Go/src/runtime/panic.go:113
//         C:/Users//Documents/gostudy/http-web/day7-error/main.go:17
//         C:/Users//Documents/gostudy/http-web/day7-error/engine/context.go:65
//         C:/Users//Documents/gostudy/http-web/day7-error/engine/recovery.go:24
//         C:/Users//Documents/gostudy/http-web/day7-error/engine/context.go:65
//         C:/Users//Documents/gostudy/http-web/day7-error/engine/logger.go:15
//         C:/Users//Documents/gostudy/http-web/day7-error/engine/context.go:65
//         C:/Users//Documents/gostudy/http-web/day7-error/engine/router.go:101
//         C:/Users//Documents/gostudy/http-web/day7-error/engine/engine.go:154
//         C:/Program Files/Go/src/net/http/server.go:2948
//         C:/Program Files/Go/src/net/http/server.go:1992
//         C:/Program Files/Go/src/runtime/asm_amd64.s:1595

// 2024/06/07 10:03:59 [500] /panic in 428.3µs

func main() {
	r := engine.Default()
	r.Get("/", func(c *engine.Context) {
		c.Plain(http.StatusOK, "Hello test\n")
	})
	// index out of range for testing Recovery()
	r.Get("/panic", func(c *engine.Context) {
		names := []string{"test"}
		c.Plain(http.StatusOK, names[100])
	})

	r.Run(":8080")
}
```