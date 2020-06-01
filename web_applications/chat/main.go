package main

import (
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"flag"
	// "os"
	// "mycode/trace"
)

// templはもう一つのテンプレートを表す
type templateHandler struct {
	once	 sync.Once
	filename string
	templ	 *template.Template
}

//　HTTPリクエストを処理 ServeHTTPという名前以外だと動かない
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join("templates",
				t.filename)))
	})
	t.templ.Execute(w, r)
}

func main(){
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス") // *string型を返すので、フラグの値は*演算子をつける必要がある
	flag.Parse() // フラグを解釈する
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"}) //　ルート
	http.Handle("/room", r)
	go r.run()
	log.Println("Webサーバーを開始。ポート:", *addr) // // *string型を返すので、フラグの値は*演算子をつける必要がある
	if err := http.ListenAndServe(":8080",nil); err != nil { // Webサーバーを開始
		log.Fatal("ListenAndServe:", err)
	}
}