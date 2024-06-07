package main

import (
	"net/http"

	"engine"
)

// $ go run .
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
