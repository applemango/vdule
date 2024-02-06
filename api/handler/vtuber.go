package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"vdule/vtuber"
)

func GetVtubers(c *gin.Context) {
	vtubers, err := vtuber.GetVtubers()
	if err != nil {
		fmt.Printf("[error] %v", err.Error())
		c.JSON(http.StatusInternalServerError, "500")
		return
	}
	c.JSON(http.StatusOK, vtubers)
}
