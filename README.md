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