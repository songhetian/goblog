package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

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

	err := r.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "请提供正确的数据")

		return
	}

	title := r.PostForm.Get("title")

	fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	fmt.Fprintf(w, "title 值: %v <br>", title)

	fmt.Fprintf(w, "r.Form 中的title值为：%v <br>", r.FormValue("title"))
	fmt.Fprintf(w, "r.PostForm 中的title值为：%v <br>", r.PostFormValue("title"))

	fmt.Fprintf(w, "r.Form 中的test值为：%v <br>", r.FormValue("test"))
	fmt.Fprintf(w, "r.PostForm 中的test值为：%v <br>", r.PostFormValue("test"))

}

func articlesCreatedHandler(w http.ResponseWriter, r *http.Request) {

	html := `<!DOCTYPE html>
				<html lang="en">
				<head>
					<title>创建文章 —— 我的技术博客</title>
				</head>
				<body>
					<form action="%s?test=data" method="post">
						<p><input type="text" name="title"></p>
						<p><textarea name="body" cols="30" rows="10"></textarea></p>
						<p><button type="submit">提交</button></p>
					</form>
				</body>
				</html>
				`
	storeURL, _ := router.Get("articles.store").URL()

	fmt.Fprintf(w, html, storeURL)
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//设置表头

		w.Header().Set("content-type", "text/html ; charset=utf-8")

		//继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 中间件
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandle).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoredHandler).Methods("POST").Name("articles.store")

	router.HandleFunc("/articles/create", articlesCreatedHandler).Methods("GET").Name("articles.create")
	//自定义404页面

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//中间件:强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
