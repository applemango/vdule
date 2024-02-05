package main

import (
	"github.com/gin-gonic/gin"
	"vdule/api/handler"
	"vdule/api/middleware"
	"vdule/youtube"
)

func main() {

	if youtube.T == nil {
		panic("youtube connection error")
	}
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(middleware.Cors())
	app.GET("/channel", handler.GetChannel)
	app.GET("/schedule", handler.GetSchedules)
	_ = app.Run(":8081")
}
