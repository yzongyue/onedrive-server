package main

import "github.com/gorilla/websocket"

type ClientHub interface {
	SendMsg(sid, msg string)
	AddClient(c *Client)
	RemoveClient(c *Client)
}

func NewClientHub() ClientHub {
	return &clientHub{clients: map[string]*Client{}}
}

type clientHub struct {
	clients map[string]*Client
}

func (hub *clientHub) SendMsg(sid, msg string) {
	client := hub.getClientBySid(sid)
	if client == nil {
		sugar.Warnf("conn not found: %s", sid)
		return
	} else {
		if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			client.Close()
			sugar.Warnf("ws send err, %s:%s", sid, msg)
		}
	}
}

func (hub *clientHub) getClientBySid(sid string) *Client {
	if v, ok := hub.clients[sid]; ok {
		return v
	}
	return nil
}

func (hub *clientHub) AddClient(c *Client) {
	// en... one conn now, err TODO
	if v, ok := hub.clients[c.sid]; ok {
		hub.RemoveClient(v)
	}
	hub.clients[c.sid] = c
	go c.wsRead()
}

func (hub *clientHub) RemoveClient(c *Client) {
	if v, ok := hub.clients[c.sid]; ok {
		if v == c {
			delete(hub.clients, c.sid)
		}
	}
	c.Close()
}

type Client struct {
	conn *websocket.Conn
	sid  string
	hub  ClientHub
}

func (c *Client) Close() {
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
	hub.RemoveClient(c)
}

func (c *Client) WriteMessage(t int, msg []byte) error {
	return c.conn.WriteMessage(t, msg)
}
