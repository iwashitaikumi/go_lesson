package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"mycode/trace"
)

type room struct {
	forward chan []byte  	 // 他のクライアントに転送するためのメッセージを保持する
	join    chan *client 	 // チャットルームに参加しようとしているクライアントのためのチャネル
	leave   chan *client 	 // チャットルームから退室しようとしているクライアントのためのチャネル
	clients map[*client]bool // 在室している全てのクライアントが保持される
	tracer trace.Tracer 	 // tracerはチャットルーム上で行われた操作のログを受け取る
}

func newRoom() *room {
	return &room {
		forward: make(chan []byte),
		join:  	 make(chan *client),
		leave:	 make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r * room) run() {
	for {
		select {
		case client := <- r.join:
			r.clients[client] = true // 参加
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <- r.leave:
			delete(r.clients, client) // 退室
			close(client.send)
			r.tracer.Trace("クライアントが退出しました")
		case msg := <- r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace(" -- クライアントに送信されました")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンアップします")
				}
			}
		}
	}
}

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, reg *http.Request) {
	socket, err := upgrader.Upgrade(w, reg, nil)
	if err != nil {
		log.Fatal("ServeHTTP",err)
		return
	}
	client := &client{
		socket: socket,
		send:	make(chan []byte, messageBufferSize),
		room:	r,
	}
	r.join <- client
	defer func() {r.leave <- client}() // クライアントの終了時にクリーンアップ
	go client.write()
	client.read()
}