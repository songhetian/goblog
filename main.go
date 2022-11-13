package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		fmt.Fprintf(w, "<h1>hello,golang2</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "<h1>请求的地址1是:"+r.URL.Path+"</h1>")
	} else {
		fmt.Fprint(w, "<h1>请求的地址不存在</h1>")
	}
}
func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3001", nil)
}
