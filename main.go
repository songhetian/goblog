package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "<h1>hello,golang2</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求的地址不存在</h1>")
	}
}

func aboutHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.URL.Path == "/about" {
		fmt.Fprintf(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
			"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
	}
}

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/", handlerFunc)
	router.HandleFunc("/about", aboutHander)

	//文章详情
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {

		// id := strings.SplitN(r.URL.Path, "/", 2)[1]
		// fmt.Fprint(w, "文章 ID："+id)

		switch r.Method {
		case "GET":
			fmt.Fprint(w, "访问文章列表")
		case "POST":
			fmt.Fprint(w, "创建新的文章")

		}
	})

	http.ListenAndServe(":3000", router)
}
