package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}
	fmt.Println(config.FormatDSN())
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = db.Ping()
	checkError(err)
}

type Article struct {
	Title, Body string
	ID          int64
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

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

	//读取对应的数据结构
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {

		tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)

		err = tmpl.Execute(w, article)
		checkError(err)
		//fmt.Fprintf(w, "读取成功，文章标题----"+article.Title)
	}

}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "访问文章列表")
}

func articlesStoredHandler(w http.ResponseWriter, r *http.Request) {

	// err := r.ParseForm()

	// if err != nil {
	// 	fmt.Fprintf(w, "请提供正确的数据")

	// 	return
	// }

	// title := r.PostForm.Get("title")

	// fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	// fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	// fmt.Fprintf(w, "title 值: %v <br>", title)

	// fmt.Fprintf(w, "r.Form 中的title值为：%v <br>", r.FormValue("title"))
	// fmt.Fprintf(w, "r.PostForm 中的title值为：%v <br>", r.PostFormValue("title"))

	// fmt.Fprintf(w, "r.Form 中的test值为：%v <br>", r.FormValue("test"))
	// fmt.Fprintf(w, "r.PostForm 中的test值为：%v <br>", r.PostFormValue("test"))

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	//验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需要介于3-40之间"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	if len(errors) == 0 {
		LastInsertId, err := saveArticleToDB(title, body)
		if LastInsertId > 0 {
			fmt.Fprint(w, "插入成功,ID为"+strconv.FormatInt(LastInsertId, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {

		storeUrl, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeUrl,
			Errors: errors,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}

	}

}

func saveArticleToDB(title string, body string) (int64, error) {
	//设置初始化

	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	//获取一个perpare声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	//例行的错误检查
	if err != nil {
		return 0, err
	}

	//在此函数运行结束后关闭此语句 防止占用sql连接
	defer stmt.Close()

	//执行请求 传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	//插入成功的话 会返回自增id
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

func articlesCreatedHandler(w http.ResponseWriter, r *http.Request) {

	storeURL, _ := router.Get("articles.store").URL()

	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}

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

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func main() {
	initDB()
	createTables()

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
