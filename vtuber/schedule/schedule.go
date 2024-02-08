package schedule

import (
	"database/sql"
	youtube2 "google.golang.org/api/youtube/v3"
	"time"
	db "vdule/utils/db/sqlite3"
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

type Row interface {
	*sql.Row | *sql.Rows
	Scan(dest ...any) error
}

func GetScheduleCore[T Row](row T) (*Schedule, error) {
	var (
		s Schedule
		d ScheduleDate
	)
	var handle string
	err := row.Scan(&s.Id, &handle, &s.Title, &s.Thumbnail, &s.IsNowOnAir, &d.Year, &d.Month, &d.Day, &d.Hour, &d.Minute)
	if err != nil {
		return nil, err
	}
	channel, in := youtube.GetRawChannelByHandleCache(handle)
	if !in {
		return nil, err
	}
	s.Channel = youtube.ChannelToTubeChannel(channel)
	s.Date = d
	return &s, nil
}

func GetSchedulesCore(query string, args ...any) ([]Schedule, error) {
	var schedules []Schedule
	rows, err := db.Conn.Query(query, args...)
	if err != nil {
		return schedules, err
	}
	for rows.Next() {
		schedule, err := GetScheduleCore(rows)
		if err != nil {
			continue
		}
		schedules = append(schedules, *schedule)
	}
	return schedules, nil
}

func GetSchedules(year, month, day int) ([]Schedule, error) {
	return GetSchedulesCore(`SELECT id, handle, title, thumbnail, is_now_on_air, live_scheduled_start_year, live_scheduled_start_month, live_scheduled_start_day, live_scheduled_start_hour, live_scheduled_start_minute FROM video WHERE live_scheduled_start_year = ? AND live_scheduled_start_month = ? AND live_scheduled_start_day = ? AND is_live = true`, year, month, day)
}

func GetChannelSchedules(year, month, day int, handle string) ([]Schedule, error) {
	return GetSchedulesCore(`SELECT id, handle, title, thumbnail, is_now_on_air, live_scheduled_start_year, live_scheduled_start_month, live_scheduled_start_day, live_scheduled_start_hour, live_scheduled_start_minute FROM video WHERE live_scheduled_start_year = ? AND live_scheduled_start_month = ? AND live_scheduled_start_day = ? AND is_live = true AND handle = ?`, year, month, day, handle)
}

func AddScheduleFromRawVideo(channel youtube2.Channel, video youtube2.Video) error {
	c := youtube.ChannelToTubeChannel(&channel)
	v := youtube.VideoToTubeVideo(&video)
	return AddScheduleFromVideo(AddScheduleProps{
		VideoId:    video.Id,
		ChannelId:  channel.Id,
		Handle:     c.Handle,
		Title:      v.Title,
		Thumbnail:  v.Thumbnail,
		IsLive:     v.IsLive,
		IsNowOnAir: v.IsLiveNowOnAir,
		Date:       v.LiveDate,
	})
}

func AddScheduleFromVideo(p AddScheduleProps) error {
	_ = DeleteSomeSchedule(p)
	_, err := db.Conn.Exec(
		`INSERT INTO video (id, channel_id, handle, title, thumbnail, is_live, is_now_on_air, live_scheduled_start_year, live_scheduled_start_month, live_scheduled_start_day, live_scheduled_start_hour, live_scheduled_start_minute ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.VideoId, p.ChannelId, p.Handle, p.Title, p.Thumbnail, p.IsLive, p.IsNowOnAir, p.Date.Year(), p.Date.Month(), p.Date.Day(), p.Date.Hour(), p.Date.Minute(),
	)
	return err
}

type AddScheduleProps struct {
	VideoId    string
	ChannelId  string
	Handle     string
	Title      string
	Thumbnail  string
	IsLive     bool
	IsNowOnAir bool
	Date       time.Time
}

func DeleteNearTimeSchedule(handle string, date time.Time) error {
	_, err := db.Conn.Exec(`DELETE FROM schedule WHERE handle = ? AND time_year = ? AND time_month = ? AND time_day = ? AND ( time_hour >= ? - 1 AND time_hour <= ? + 1 ) `, handle, date.Year(), date.Month(), date.Day(), date.Hour(), date.Hour())
	return err
}

func DeleteSomeSchedule(p AddScheduleProps) error {
	_, err := db.Conn.Exec(`DELETE FROM video WHERE id = ?`, p.VideoId)
	if err != nil {
		return err
	}
	err = DeleteNearTimeSchedule(p.Handle, p.Date)
	return err
}

func AddSchedule(p AddScheduleProps) error {
	_ = DeleteNearTimeSchedule(p.Handle, p.Date)
	_, err := db.Conn.Exec(
		`INSERT INTO schedule (channel_id, handle, title, thumbnail, time_year, time_month, time_day, time_hour ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ChannelId, p.Handle, p.Title, p.Thumbnail, p.Date.Year(), p.Date.Month(), p.Date.Day(), p.Date.Hour(),
	)
	return err
}
