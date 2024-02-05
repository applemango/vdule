package schedule

import (
	"fmt"
	"strings"
	"time"
	db "vdule/db/sqlite3"
	"vdule/utils/http"
	"vdule/youtube"
)

// https://schedule.hololive.tv/api/list

type HololiveScheduleApiResult struct {
	DateGroupList []struct {
		DisplayDate string `json:"displayDate"`
		Datetime    string `json:"datetime"`
		VideoList   []struct {
			DisplayDate  string `json:"displayDate"`
			Datetime     string `json:"datetime"`
			IsLive       bool   `json:"isLive"`
			PlatformType int    `json:"platformType"`
			URL          string `json:"url"`
			Thumbnail    string `json:"thumbnail"`
			Title        string `json:"title"`
			Name         string `json:"name"`
			Talent       struct {
				Name         string `json:"name"`
				IconImageURL string `json:"iconImageUrl"`
			} `json:"talent"`
			CollaboTalents []interface{} `json:"collaboTalents"`
		} `json:"videoList"`
	} `json:"dateGroupList"`
}

func GetHololiveRawSchedule() (*HololiveScheduleApiResult, error) {
	var data *HololiveScheduleApiResult
	if err := http.HttpGetStructCacheExp("https://schedule.hololive.tv/api/list", &data, time.Hour); err != nil {
		return nil, err
	}
	return data, nil
}

func ParseHololiveApiDate(datetime string) time.Time {
	date, err := time.Parse("2006/01/02 15:04:05", datetime)
	if err != nil {
		panic(err.Error())
	}
	return date
}

func ParseYoutubeVideoId(url string) string {
	id := strings.Split(url, "=")[1]
	return id
}

var HololiveChannelId = map[string]string{
	"ときのそら":  "TokinoSora",
	"博衣こより":  "HakuiKoyori",
	"湊あくあ":   "MinatoAqua",
	"白上フブキ":  "ShirakamiFubuki",
	"さくらみこ":  "SakuraMiko",
	"猫又おかゆ":  "NekomataOkayu",
	"常闇トワ":   "TokoyamiTowa",
	"大空スバル":  "OozoraSubaru",
	"百鬼あやめ":  "NakiriAyame",
	"兎田ぺこら":  "usadapekora",
	"白銀ノエル":  "ShiroganeNoel",
	"角巻わため":  "TsunomakiWatame",
	"天音かなた":  "AmaneKanata",
	"Gura":   "GawrGura",
	"沙花叉クロヱ": "SakamataChloe",
	"ロボ子さん":  "Robocosan",
	"雪花ラミィ":  "YukihanaLamy",
}

func RegisterHololiveSchedule() error {
	res, err := GetHololiveRawSchedule()
	if err != nil {
		return err
	}
	for _, group := range res.DateGroupList {
		for _, video := range group.VideoList {
			handle, in := HololiveChannelId[video.Talent.Name]
			if !in {
				fmt.Printf("[info] not in hololive channel handles map: %v\n", video.Talent.Name)
				continue
			}
			if video.PlatformType != 1 {
				continue
			}
			channel, err := youtube.T.GetChannelByHandle(handle)
			if err != nil {
				return err
			}
			id := ParseYoutubeVideoId(video.URL)
			date := ParseHololiveApiDate(video.Datetime)
			_, _ = db.Conn.Exec(`DELETE FROM video WHERE id = ?`, id)
			_, err = db.Conn.Exec(`INSERT INTO video (
			   	id,
			   	channel_id,
                handle,
                title,
                thumbnail,
			   	is_live,
                is_now_on_air,
			   	live_scheduled_start_year,
			   	live_scheduled_start_month,
			   	live_scheduled_start_day,
			   	live_scheduled_start_hour,
            	live_scheduled_start_minute
			) VALUES (?, ?, ?, ?, ?, true, ?, ?, ?, ?, ?, ?)`, id, channel.Id, youtube.ParseChannelHandle(handle), video.Title, video.Thumbnail, video.IsLive, date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
