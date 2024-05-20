package engine

import (
	"net/http"
	"strings"
)

// Improvement: support dynamic routing
// introducr router root map[string]Node to build one trie tree per REQ Method, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// split pattern for routing path match, only one * is allowed
func parsePattern(pattern string) []string {
	pList := make([]string, 0)
	for _, e := range strings.Split(pattern, "/") {
		if e != "" {
			pList = append(pList, e)
			if e[0] == '*' {
				break
			}
		}
	}
	return pList
}

// add routing rules to router table
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	pList := parsePattern(pattern)
	// fmt.Println(pList)
	// construct root
	_, ok := r.roots[method]
	if !ok {
		// init trie method tree
		r.roots[method] = &node{}
	}
	// insert height 0 starting from top
	r.roots[method].insert(pattern, pList, 0)
	// handler key
	hKey := method + "-" + pattern
	r.handlers[hKey] = handler
}

// search route table for path and return node and updated params map to be used in context
func (r *router) searchRoute(method string, path string) (*node, map[string]string) {
	pathList := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	node := root.search(pathList, 0)
	if node != nil {
		patternList := parsePattern(node.pattern)
		for i, e := range patternList {
			// patternList /users/:userId/posts/:postId
			// pathList is ["users", "123", "posts", "456"]
			// params {"userId": "123", "postId": "456"}
			if e[0] == ':' {
				params[e[1:]] = pathList[i]
			}
			// {pattern: "/files/*filepath"}
			// pathList := []string{"files", "images", "2024", "may", "photo.jpg"}
			// params map will then contain the key "filepath" with the value "images/2024/may/photo.jpg"
			if e[0] == '*' && len(e) > 1 {
				params[e[1:]] = strings.Join(pathList[i:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func (r *router) handleRoute(c *Context) {
	node, params := r.searchRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		// note key is pattern, not path
		key := c.Method + "-" + node.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.Plain(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// after appended the router handler itself, we start the middleware chain execution
	c.Next()
}

// get all route entries of given method, i.e. return
// all leaf nodes (with pattern defined)
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.getPatternNodes(&nodes)
	return nodes
}
