package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type WebHookQuery struct {
	ValidationToken string `form:"validationToken"`
}

type WebHookMsg struct {
	Value []MsgItem `json:"value"`
}

type MsgItem struct {
	ChangeType                     string      `json:"changeType"`
	ClientState                    string      `json:"clientState"`
	Resource                       string      `json:"resource"`
	ResourceData                   interface{} `json:"resourceData"`
	SubscriptionExpirationDateTime string      `json:"subscriptionExpirationDateTime"`
	SubscriptionID                 string      `json:"subscriptionId"`
	TenantID                       string      `json:"tenantId"`
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
				go hub.SendMsg(v.SubscriptionID, string(m))
			}
		}
	}
}
