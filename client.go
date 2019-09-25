package main

import "github.com/gorilla/websocket"

var clients = map[string]*Client{}

func getClientBySid(sid string) *Client {
	if v, ok := clients[sid]; ok {
		return v
	}
	return nil
}

type Client struct {
	conn *websocket.Conn
	sid  string
}

func (c *Client) Run() {
	// en... one conn now, err TODO
	if _, ok := clients[c.sid]; ok {
		clients[c.sid].Close()
	}
	clients[c.sid] = c
	go c.wsRead()
}

func (c *Client) Close() {
	if v, ok := clients[c.sid]; ok {
		if v == c {
			delete(clients, c.sid)
		}
	}
	_ = c.conn.Close()
}

func (c *Client) wsRead() {
	c.conn.SetReadLimit(16) // TODO
	for {
		t, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if t == websocket.PingMessage {
			_ = c.conn.WriteMessage(websocket.PongMessage, message)
		}
	}
	c.Close()
}

func (c *Client) WriteMessage(t int, msg []byte) error {
	return c.conn.WriteMessage(t, msg)
}
