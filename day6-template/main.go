package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"engine"
)

// ---engine/
// ---static/
//    |---css/
//         |---geektutu.css
//    |---file1.txt
// ---templates/
//    |---arr.tmpl
//    |---css.tmpl
//    |---custom_func.tmpl
// ---main.go

/*
(1) render array
$ curl http://localhost:9999/date
<html>
<body>
    <p>hello, gee</p>
    <p>Date: 2019-08-17</p>
</body>
</html>
*/

/*
(2) custom render function
$ curl http://localhost:9999/students
<html>
<body>
    <p>hello, gee</p>
    <p>0: Geektutu is 20 years old</p>
    <p>1: Jack is 22 years old</p>
</body>
</html>
*/

/*
(3) serve static files via assets/
$ curl http://localhost:9999/assets/css/geektutu.css
p {
    color: orange;
    font-weight: 700;
    font-size: 20px;
}
$ curl http://localhost:9999/assets/css/geektutu1.css -v
*   Trying 127.0.0.1:9999...
* Connected to localhost (127.0.0.1) port 9999 (#0)
> GET /assets/css/geektutu1.css HTTP/1.1
> Host: localhost:9999
> User-Agent: curl/7.80.0
> Accept:
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 404 Not Found
< Date: Wed, 22 May 2024 07:08:59 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
*/

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
}
