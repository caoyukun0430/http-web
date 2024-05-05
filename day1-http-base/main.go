package main

import (
	"fmt"
	"net/http"

	"engine"
)

// func main() {
// 	http.HandleFunc("/", rootHandler)
// 	http.HandleFunc("/hello", helloHandler)
// 	fmt.Printf("http server running ...\n")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// // handler echoes r.URL.Path
// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "URL Path is %q\n", r.URL.Path)
// }

// // handler echoes r.URL.Header
// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	for k, v := range r.Header {
// 		fmt.Fprintf(w, "Header[%q]: %q\n", k, v)
// 	}
// }

// Improvement 1: instead of calling individual HandleFunc() using individual handlers,
// we implement interface http.handler by defining serveHTTP() and redirect all routing logic
// to our own Engine handler, overall doesnt make much difference between a single Handler
// and multi-HandleFunc, important is we abstract the Engine type and implements routing inside

// type Engine struct{}

// func (engine Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		fmt.Fprintf(w, "URL Path is %q\n", r.URL.Path)
// 	case "/hello":
// 		for k, v := range r.Header {
// 			fmt.Fprintf(w, "Header[%q]: %q\n", k, v)
// 		}
// 	default:
// 		fmt.Fprintf(w, "default URL Path is %q\n", r.URL.Path)
// 	}
// }

// func main() {
// 	engine := new(Engine)
// 	fmt.Printf("http server running ...\n")
// 	log.Fatal(http.ListenAndServe(":8080", engine))
// }

// Improvment 2: In order to diff GET/POST, we extend our routing table by string-handler map
// map[string]HandlerFunc where key is METHOD-PATH, and val is HandlerFunc
// key idea is Engine still handles all routing and inside ServeHTTP(), it disributes works
// to individual HandlerFunc in the next level, which is ordinary func to implements Handler
// interface simply
func main() {
	server := engine.New()
	// defines pattern and handler as HandlerFunc method
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
