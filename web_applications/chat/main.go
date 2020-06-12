package main

import (
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
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
	data := map[string]interface{} {
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main(){
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス") // *string型を返すので、フラグの値は*演算子をつける必要がある
	flag.Parse() // フラグを解釈する
	gomniauth.SetSecurityKey("セキュリティー") // gomniauthのセットアップ
	gomniauth.WithProviders (
		facebook.New("958358529337-029gmndul8jgn86dqfjaaiiu4h74gjlc.apps.googleusercontent.com","Gg9eBuEUTt53dzO4_MVG7Llc","http://localhost:8080/auth/callback/facebook"),
		github.New("958358529337-029gmndul8jgn86dqfjaaiiu4h74gjlc.apps.googleusercontent.com","Gg9eBuEUTt53dzO4_MVG7Llc","http://localhost:8080/auth/callback/github"),
		google.New("958358529337-029gmndul8jgn86dqfjaaiiu4h74gjlc.apps.googleusercontent.com","Gg9eBuEUTt53dzO4_MVG7Llc","http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	http.Handle("/assets/", http.StripPrefix("/assets",http.FileServer(http.Dir("/assets"))))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"})) //　ルート
	http.Handle("/login",&templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/logout",func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie {
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	go r.run()
	log.Println("Webサーバーを開始。ポート:", *addr) // // *string型を返すので、フラグの値は*演算子をつける必要がある
	if err := http.ListenAndServe(":8080",nil); err != nil { // Webサーバーを開始
		log.Fatal("ListenAndServe:", err)
	}
}