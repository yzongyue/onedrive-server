package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleWS(c *gin.Context) {
	sid := c.Param("sid")
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		sugar.Errorf("ws upgrade err %s: %s", sid, err)
		return
	}

	client := &Client{conn: conn, sid: sid, hub: hub}
	hub.AddClient(client)
}
