package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"vdule/vtuber"
	"vdule/vtuber/schedule"
	"vdule/vtuber/youtube"
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

func GetVtubersDetail(c *gin.Context) {
	if handle, in := c.GetQuery("handle"); in {
		var (
			v                   *vtuber.Vtuber
			organizationMembers []vtuber.Vtuber
		)
		v, err := vtuber.GetVtuberByHandle(handle)
		if err != nil {
			channel, _ := youtube.T.GetChannelByHandle(handle)
			v = &vtuber.Vtuber{
				Channel: *channel,
			}
		}
		if v.Organization != nil {
			organizationMembers, err = vtuber.GetVtubersByOrganization(v.Organization.Id)
		}
		now := time.Now()
		year, month, day := now.Year(), int(now.Month()), now.Day()
		schedules, err := schedule.GetChannelSchedules(year, month, day, v.Channel.Handle)
		c.JSON(http.StatusOK, map[string]any{
			"vtuber":               &v,
			"organization_members": organizationMembers,
			"schedule":             schedules,
		})
		return
	}
	c.JSON(http.StatusNotFound, "Not Found.")
}
