package main

import (
	"github.com/stretchr/objx"
	"net/http"
	"strings"
	"log"
	"fmt"
	"github.com/stretchr/gomniauth"
	"crypto/md5"
	"io"
)

type authHundler struct {
	next http.Handler
}

func (h *authHundler) ServeHTTP(w http.ResponseWriter, r * http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie { // 未認証
		w.Header().Set("Location","/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil { // 何らかのエラーが発生
		panic(err.Error())
	} else { // 成功
		h.next.ServeHTTP(w,r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHundler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) { // /auth/{action}/{provider}
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダの取得に失敗しました:",provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました:", provider, "-", err)
		}
		w.Header().Set("Location",loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダの取得に失敗しました:",provider, "-", err)
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした", provider, "-", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザの取得に失敗しました",provider, "-", err)
		}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Name()))
		userID := fmt.Sprintf("%x", m.Sum(nil))
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     userID,
			"name":		  user.Name(),
			"avatar_url": user.AvatarURL(),
			"email":	  user.Email(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}