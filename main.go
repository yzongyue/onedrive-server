package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
	hub   ClientHub
)

func init() {
	sugar = zap.NewExample().Sugar()
	hub = NewClientHub()
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/webhook", HandleWebHook)
	r.GET("/ws/:sid/:rnd", HandleWS)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	return r
}

func main() {
	defer sugar.Sync()
	r := setupRouter()
	r.Run("localhost:6500")
}
