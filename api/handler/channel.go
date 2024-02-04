package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vdule/youtube"
)

func GetChannel(c *gin.Context) {
	if handle, in := c.GetQuery("handle"); in {
		if channel, err := youtube.T.GetChannelByHandle(handle); err == nil {
			c.JSON(http.StatusOK, channel)
			return
		}
	}
	c.JSON(http.StatusNotFound, "Not found.")
}
