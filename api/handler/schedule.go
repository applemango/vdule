package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"vdule/utils/db/sqlite3"
	"vdule/vtuber/youtube"
)

type ScheduleDate struct {
	Year   int `json:"year,omitempty"`
	Month  int `json:"month,omitempty"`
	Day    int `json:"day,omitempty"`
	Hour   int `json:"hour,omitempty"`
	Minute int `json:"minute,omitempty"`
}

type Schedule struct {
	Id         string              `json:"id,omitempty"`
	Date       ScheduleDate        `json:"date"`
	Channel    youtube.TubeChannel `json:"channel"`
	IsNowOnAir bool                `json:"is_now_on_air,omitempty"`
	Title      string              `json:"title,omitempty"`
	Thumbnail  string              `json:"thumbnail,omitempty"`
}

func GetSchedules(c *gin.Context) {
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()
	var handle *string
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
	var rows *sql.Rows
	var err error
	if handle != nil {
		rows, err = db.Conn.Query(`SELECT
    		id,
    		handle,
    		title,
    		thumbnail,
    		is_now_on_air,
    		live_scheduled_start_year,
    		live_scheduled_start_month,
    		live_scheduled_start_day,
    		live_scheduled_start_hour,
    		live_scheduled_start_minute
		FROM video WHERE live_scheduled_start_year = ? AND live_scheduled_start_month = ? AND live_scheduled_start_day = ? AND is_live = true AND handle = ?`, year, month, day, handle)
	} else {
		rows, err = db.Conn.Query(`SELECT
    		id,
    		handle,
    		title,
    		thumbnail,
    		is_now_on_air,
    		live_scheduled_start_year,
    		live_scheduled_start_month,
    		live_scheduled_start_day,
    		live_scheduled_start_hour,
    		live_scheduled_start_minute
		FROM video WHERE live_scheduled_start_year = ? AND live_scheduled_start_month = ? AND live_scheduled_start_day = ? AND is_live = true`, year, month, day)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, "500")
		return
	}
	var res []Schedule
	for rows.Next() {
		var s Schedule
		var d ScheduleDate
		var handle string
		err := rows.Scan(&s.Id, &handle, &s.Title, &s.Thumbnail, &s.IsNowOnAir, &d.Year, &d.Month, &d.Day, &d.Hour, &d.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "500")
			return
		}
		channel, in := youtube.GetRawChannelByHandleCache(handle)
		if !in {
			continue
		}
		s.Channel = youtube.ChannelToTubeChannel(channel)
		s.Date = d
		res = append(res, s)
	}
	c.JSON(http.StatusOK, map[string]any{
		"videos": res,
	})
}
