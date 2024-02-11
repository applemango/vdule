package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"vdule/api/handler"
	"vdule/api/middleware"
	"vdule/utils/cache"
	y "vdule/vtuber/youtube"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("failed load config")
	}
	y.ResetYoutube()
	if y.T == nil {
		panic("youtube connection error")
	}
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.Use(middleware.Cors())
	app.GET("/channel", handler.GetChannel)
	app.GET("/schedule", cache.ResponseMiddleware(handler.GetSchedules))
	app.GET("/vtubers", handler.GetVtubers)
	app.GET("/vtuber", handler.GetVtubersDetail)
	_ = app.Run(":8081")
}
