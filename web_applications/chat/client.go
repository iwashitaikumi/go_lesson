package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// クライアントはチャットを行っている1人のユーザを表す
type client struct {
	socket   *websocket.Conn        // このクライアントのためのwebsoclet
	send     chan *message          // メッセージが送られるチャンネル
	room     *room                  // このクライアントが参加しているチャットルーム
	userData map[string]interface{} // ユーザに関する情報
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			now := time.Now()
			jst := time.FixedZone("Asia/Tokyo", 9*60*60)
			nowUTC := now.UTC()
			nowJST := nowUTC.In(jst)
			msg.When = nowJST                    
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}