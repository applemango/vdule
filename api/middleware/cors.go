package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1:3001",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
		},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"Account",
		},
		MaxAge: 0,
	})
}
