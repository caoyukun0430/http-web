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