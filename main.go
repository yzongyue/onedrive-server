package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
)

type WebHookQuery struct {
	ValidationToken string `form:"validationToken"`
}

type WebHookMsg struct {
	Value []struct {
		ChangeType                     string      `json:"changeType"`
		ClientState                    string      `json:"clientState"`
		Resource                       string      `json:"resource"`
		ResourceData                   interface{} `json:"resourceData"`
		SubscriptionExpirationDateTime string      `json:"subscriptionExpirationDateTime"`
		SubscriptionID                 string      `json:"subscriptionId"`
		TenantID                       string      `json:"tenantId"`
	} `json:"value"`
}

func HandleWebHook(c *gin.Context) {
	var query WebHookQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.String(400, "query err")
		return
	}
	if query.ValidationToken != "" {
		c.String(200, query.ValidationToken)
		return
	} else {
		var msg WebHookMsg
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.String(400, "json err")
			sugar.Warn("parse webhook json err")
			return
		} else {
			for _, v := range msg.Value {
				m, _ := json.Marshal(v)
				go SendMsg(v.SubscriptionID, string(m))
			}
		}
	}
}

func SendMsg(sid string, msg string) {
	client := getClientBySid(sid)
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

func HandleWS(c *gin.Context) {
	sid := c.Param("sid")
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		sugar.Errorf("ws upgrade err %s: %s", sid, err)
		return
	}

	client := Client{conn: conn, sid: sid}
	client.Run()
}

func main() {
	sugar = zap.NewExample().Sugar()
	defer sugar.Sync()

	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/webhook", HandleWebHook)
	r.GET("/ws/:sid/:rnd", HandleWS)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	r.Run("localhost:6500")
}
