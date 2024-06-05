package engine

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/tom", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/alice", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/a"), []string{"p", "a"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

// concrete path match has priority, if present, will be selected
func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, params := r.searchRoute("GET", "/hello/alice")

	fmt.Printf("matched path: %s\n", n.pattern)

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/alice" {
		t.Fatal("should match /hello/alice")
	}

	n, params = r.searchRoute("GET", "/hello/tom")
	// concrete path match has priority
	fmt.Printf("matched path: %s\n", n.pattern)

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/tom" {
		t.Fatal("should match /hello/tom")
	}

	n, params = r.searchRoute("GET", "/hello/bob")

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, params["name"])

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if params["name"] != "bob" {
		t.Fatal("name should be equal to 'bob'")
	}

	n, params = r.searchRoute("GET", "/assets/image/2024/1.jpg")

	fmt.Printf("matched path: %s, params['filepath']: %s\n", n.pattern, params["filepath"])

	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if n.pattern != "/assets/*filepath" {
		t.Fatal("should match /assets/*filepath")
	}

	if params["filepath"] != "image/2024/1.jpg" {
		t.Fatal("name should be equal to 'image/2024/1.jpg'")
	}

}

func TestGetRoutes(t *testing.T) {
	// check no overwritten between /:name and /concrete_name
	r := newTestRouter()
	nodes := r.getRoutes("GET")
	for i, n := range nodes {
		fmt.Println(i+1, n)
	}

	if len(nodes) != 6 {
		t.Fatal("the number of routes shoule be 6")
	}
}
