package main

import "fmt"

// ここにfunc sayHello() を定義すると、間に合う。


func main() {
    go sayHello()
    fmt.Println("say hello from main goroutine")
}

func sayHello() { 
    fmt.Println("say hello")
}