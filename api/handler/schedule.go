package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"vdule/vtuber/schedule"
	"vdule/vtuber/youtube"
)

func GetSchedules(c *gin.Context) {
	var (
		handle *string
	)
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()
	if y, in := c.GetQuery("year"); in {
		year, _ = strconv.Atoi(y)
	}
	if m, in := c.GetQuery("month"); in {
		month, _ = strconv.Atoi(m)
	}
	if d, in := c.GetQuery("day"); in {
		day, _ = strconv.Atoi(d)
	}
	if h, in := c.GetQuery("handle"); in {
		p := youtube.ParseChannelHandle(h)
		handle = &p
	}
	if handle != nil {
		if schedules, err := schedule.GetChannelSchedules(year, month, day, *handle); err == nil {
			c.JSON(http.StatusOK, map[string]any{
				"videos": schedules,
			})
			return
		}
	}
	if schedules, err := schedule.GetSchedules(year, month, day); err == nil {
		c.JSON(http.StatusOK, map[string]any{
			"videos": schedules,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, "500")
	return
}
