package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "text/html ; charset=utf-8")
	fmt.Fprintf(w, "<h1>Hello 欢迎到来 gobolg!</h1>\n")
}

func aboutHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html ; charset=utf-8")
	fmt.Fprintf(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "文章 ID:"+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "访问文章列表")
}

func articlesStoredHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "创建新文章")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandle).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoredHandler).Methods("POST").Name("articles.store")

	//自定义404页面

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//通过命名路由获取URL示例
	homeURL, _ := router.Get("home").URL()

	fmt.Println("homeURL:", homeURL)

	articleURL, _ := router.Get("articles.show").URL("id", "23")

	fmt.Println("articleURL:", articleURL)

	http.ListenAndServe(":3000", router)
}
