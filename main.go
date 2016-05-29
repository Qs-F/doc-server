package main

import (
	"flag"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const (
	tmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta viewport="width=device-width">
	<meta charset="UTF-8">
	<title>{{.Title}}</title>
	<link rel="stylesheet" href="/style.css">
</head>
<body>
{{.Content}}
</body>
</html>
`
)

func main() {
	port := flag.Int("p", 8080, "setting application port number. default is 8080")
	flag.Parse()
	type Page struct {
		Title   string
		Content string
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wd, _ := os.Getwd()
		if v, err := os.Stat(wd + "/" + r.URL.Path[1:]); err == nil && v.Mode().IsRegular() && strings.HasSuffix(r.URL.Path[1:], ".md") { // check whether requested file exist and it is md formatted file. DONOT change the order of if checks.
			f, err := ioutil.ReadFile(wd + "/" + r.URL.Path[1:])
			if err != nil {
				http.NotFound(w, r)
				return
			}
			b := blackfriday.MarkdownCommon(f)
			tmplen, err := template.New("tmpl").Parse(tmpl)
			p := Page{
				Title:   r.URL.Path[1:],
				Content: string(b),
			}
			err = tmplen.Execute(w, p)
			if err != nil {
				http.NotFound(w, r)
			}
		} else {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	})
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), nil)
}
