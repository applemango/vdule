package main

import (
	"github.com/gin-gonic/gin"
	"vdule/api/handler"
	"vdule/api/middleware"
	youtube2 "vdule/vtuber/youtube"
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
	app.GET("/vtubers", handler.GetVtubers)
	app.GET("/vtuber", handler.GetVtubersDetail)
	_ = app.Run(":8081")
}
