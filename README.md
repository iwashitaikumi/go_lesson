# フレームワーク無し

## サンプルコード
### 主要パッケージ
* net/http$\cdots$http接続する際に必要な機能がある標準パッケージ

### コード詳細
> これはlocalhost:8080/pingにアクセスすると、pongと返してくれるアプリです。

```go
package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/ping", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
```

* L8：http.HandleFuncで、webページのルーティングを登録
	* 今回だと/pingにアクセスした時に、handler関数を呼び出すという処理を登録している
* L9：http.ListenAndServeで、どのポートで公開するかを指定
	* 今回だと他のポートと被らないように、8080ポートで公開するようにしている

* L12
	* http.ResponseWriter：レスポンスに関する構造体→情報を返す時に利用
	* http.Request：リクエストに関する構造体→リクエストに含まれている情報を取得する時に利用
* L13：w.WriteHeader(http.StatusOK)は レスポンスに倒するステータスコードを記述
	* 200でもいいが、今回は可読性工場のために(http.StatusOK)を利用
* L14：w.Write([]byte(“pong”))で、pongという文字列をバイト列に変換して書き込んでいる

## WebAPI
### 主要パッケージ
* encoding/json$\cdots$jsonを扱うための標準パッケージ

### コード詳細
> http://localhost:8080/apiにアクセスすると、user情報がjson形式で出力されているはずです。

```go
package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api", ApiHandler)

	http.ListenAndServe(":8080", nil)
}

type User struct {
	Name string
	Age  int
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		Name: "example",
		Age:  20,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
```

* L24：w.Header().Set(“Content-Type”, “application/json”)でレスポンスのデータタイプがjson形式だと明示
* L26：http.ResponseWriterをnewEncoderの引数とする事で、レスポンスにデータを書き込む準備をし、Encode(user)でuser構造体にマッピングしている。

## Webページ
### 主要パッケージ
* html/template$\cdots$静的ファイルを生成したり、処理を加える時に使用する標準パッケージ

###  コード詳細
* http://localhost:8080にアクセスすると、user情報が表示される

#### ディレクトリ構造
```
root/
┣ main.go
┗ views/
	┗ index.tmpl
```

#### main.go

```go（
package main

import (
	"html/template"
	"net/http"fa
)

func main() {
	http.HandleFunc("/", StaticHandler)
	http.ListenAndServe(":8080", nil)
}

type User struct {
	Name string
	Age  int
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		Name: "example",
		Age:  20,
	}
	tmpl := template.Must(template.ParseFiles("./views/index.tmpl"))
	tmpl.Execute(w, user)
}
```
* L23：template.Must(template.ParseFiles(“./views/index.tmpl”))で”./views/index.tmpl”ファイルを呼び出している
* L24：tmpl.Execute(w, user)でuser構造体を静的ページに埋め込んで、レスポンスしている

#### viws/index.teml

```html
<!--- views/index.tmpl -->
<!DOCTYPE html>
<html>
<body>
    <p>Name={{.Name}}</p>
    <p>Age={{.Age}}</p>
</body>
</html>
```

* L5~：{{.Name}}でUser構造体のNameプロパティを参照しています。

## 参考ページ
[https://ryomak.info/2019/post-290/](https://ryomak.info/2019/post-290/)
