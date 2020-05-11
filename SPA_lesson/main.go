package main

import (
    "fmt"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB上のテーブル、カラムと構造体との関連付けが自動的に行われる
type Product struct {
    ID          int    `gorm:"primary_key;not null"`
    ProductName string `gorm:"type:varchar(200);not null"`
    Memo        string `gorm:"type:varchar(400)"`
    Status      string `gorm:"type:char(2);not null"`
}

func getGormConnect() *gorm.DB {
    DBMS := "mysql"
    USER := "root"
    PASS := "root"
    PROTOCOL := "tcp(localhost:3306)"
    DBNAME := "Shopping"
    CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME
    db, err := gorm.Open(DBMS, CONNECT)

    if err != nil {
        panic(err.Error())
    }

    // DBエンジンを「InnoDB」に設定
    db.Set("gorm:table_options", "ENGINE=InnoDB")

    // 詳細なログを表示
    db.LogMode(true)

    // 登録するテーブル名を単数形にする（デフォルトは複数形）
    db.SingularTable(true)

    // マイグレーション（テーブルが無い時は自動生成）
    db.AutoMigrate(&Product{})

    fmt.Println("db connected: ", &db)
    return db
}

// 商品テーブルにレコードを追加
func insertProduct(registerProduct *Product) {
    db := getGormConnect()

    // insert文
    db.Create(&registerProduct)
    defer db.Close()
}

// 商品テーブルのレコードを全件取得
func findAllProduct() []Product {
    db := getGormConnect()
    var products []Product

    // select文
    db.Order("ID asc").Find(&products)
    defer db.Close()
    return products
}

func main() {
    // Productテーブルにデータを運ぶための構造体を初期化
    var product = Product{
        ProductName: "テスト商品",
        Memo:        "テスト商品です",
        Status:      "01",
    }

    // 構造体のポインタを渡す
    insertProduct(&product)

    // Productテーブルのレコードを全件取得する
    resultProducts := findAllProduct()

    // Productテーブルのレコードを全件表示する
    for i := range resultProducts {
        fmt.Printf("index: %d, 商品ID: %d, 商品名: %s, メモ: %s, ステータス: %s\n",
            i, resultProducts[i].ID, resultProducts[i].ProductName, resultProducts[i].Memo, resultProducts[i].Status)
    }
}