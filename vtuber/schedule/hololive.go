package schedule

import (
	"fmt"
	"strings"
	"time"
	"vdule/utils/http"
	youtube2 "vdule/vtuber/youtube"
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
	/*
	 * https://schedule.hololive.tv/api/list
	 * このapiの更新頻度は15分に一度
	 */
	if err := http.HttpGetStructCacheExp("https://schedule.hololive.tv/api/list", &data, time.Minute*15); err != nil {
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
	"ときのそら":    "TokinoSora",
	"博衣こより":    "HakuiKoyori",
	"湊あくあ":     "MinatoAqua",
	"白上フブキ":    "ShirakamiFubuki",
	"さくらみこ":    "SakuraMiko",
	"猫又おかゆ":    "NekomataOkayu",
	"常闇トワ":     "TokoyamiTowa",
	"大空スバル":    "OozoraSubaru",
	"百鬼あやめ":    "NakiriAyame",
	"兎田ぺこら":    "usadapekora",
	"白銀ノエル":    "ShiroganeNoel",
	"角巻わため":    "TsunomakiWatame",
	"天音かなた":    "AmaneKanata",
	"Gura":     "GawrGura",
	"沙花叉クロヱ":   "SakamataChloe",
	"ロボ子さん":    "Robocosan",
	"雪花ラミィ":    "YukihanaLamy",
	"姫森ルーナ":    "HimemoriLuna",
	"星街すいせい":   "HoshimachiSuisei",
	"夏色まつり":    "NatsuiroMatsuri",
	"紫咲シオン":    "MurasakiShion",
	"FUWAMOCO": "FUWAMOCOch",
	"戌神ころね":    "InugamiKorone",
	"大神ミオ":     "OokamiMio",
	"ラプラス":     "LaplusDarknesss",
	"宝鐘マリン":    "HoushouMarine",
	"赤井はあと":    "AkaiHaato",
	"尾丸ポルカ":    "OmaruPolka",
	"轟はじめ":     "TodorokiHajime",
	"IRyS":     "IRyS",
	"Amelia":   "WatsonAmelia",
	"AZKi":     "AZKi",
	"桃鈴ねね":     "MomosuzuNene",
	"一条莉々華":    "IchijouRirika",
	"儒烏風亭らでん":  "JuufuuteiRaden",
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
			channel, err := youtube2.T.GetChannelByHandle(handle)
			if err != nil {
				return err
			}
			id := ParseYoutubeVideoId(video.URL)
			date := ParseHololiveApiDate(video.Datetime)
			err = AddScheduleFromVideo(AddScheduleProps{
				VideoId:    id,
				ChannelId:  channel.Id,
				Handle:     channel.Handle,
				Title:      video.Title,
				Thumbnail:  video.Thumbnail,
				IsLive:     true,
				IsNowOnAir: video.IsLive,
				Date:       date,
			})
			if err != nil {
				return err
			}
			fmt.Printf("[info] add %v's video: %v\n", video.Talent.Name, video.Name)
		}
	}
	return nil
}
