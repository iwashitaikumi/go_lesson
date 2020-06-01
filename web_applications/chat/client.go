package main

import (
	"github.com/gorilla/websocket"

)

// クライアントはチャットを行っている1人のユーザを表す
type client struct {
	socket *websocket.Conn // このクライアントのためのwebsoclet
	send   chan []byte     // メッセージが送られるチャンネル
	room   *room           // このクライアントが参加しているチャットルーム
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}