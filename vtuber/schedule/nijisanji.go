package schedule

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/exp/slices"
	"net/http"
	"strconv"
	"strings"
	"time"
	youtube2 "vdule/vtuber/youtube"
)

// https://wikiwiki.jp/nijisanji/%E9%85%8D%E4%BF%A1%E4%BA%88%E5%AE%9A%E3%83%AA%E3%82%B9%E3%83%88

func RequestNijisanjiWikiRawSchedule() (*goquery.Document, error) {
	res, err := http.Get("https://wikiwiki.jp/nijisanji/%E9%85%8D%E4%BF%A1%E4%BA%88%E5%AE%9A%E3%83%AA%E3%82%B9%E3%83%88")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func FixNijisanjiWikiRawDateNumString(str string) string {
	nums := []uint8{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	p := str
	for len(p) != 0 && !slices.Contains(nums, p[0]) {
		p = p[1:]
	}
	if len(strings.Split(p, "or")) > 1 {
		p = strings.Split(p, "or")[0]
	}
	if len(p) == 0 {
		return "00"
	}
	if len(p) == 1 {
		return fmt.Sprintf("0%v", str)
	}
	return p
}

func ParseNijisanjiWikiRawDate(year, month, day int, rawDate string) (time.Time, error) {
	s := strings.Split(rawDate[:len(rawDate)-1], "時")
	hour, minute := s[0], strings.Split(s[1], "\xe5\x88")[0]
	dateString := fmt.Sprintf("%v/%v/%v %v:%v:00", year, FixNijisanjiWikiRawDateNumString(strconv.Itoa(month)), FixNijisanjiWikiRawDateNumString(strconv.Itoa(day)), FixNijisanjiWikiRawDateNumString(hour), FixNijisanjiWikiRawDateNumString(minute))
	date, err := time.Parse("2006/01/02 15:04:05", dateString)
	if err != nil {
		return time.Now(), err
	}
	return date, nil
}

const (
	PLATFORM_YOUTUBE = iota
	PLATFORM_TWITCH
	PLATFORM_OTHER
)

func GetNijisanjiLivePlatform(title string) int8 {
	if strings.Contains(title, "at:YouTube") {
		return PLATFORM_YOUTUBE
	}
	if strings.Contains(title, "at:Twitch") {
		return PLATFORM_TWITCH
	}
	return PLATFORM_OTHER
}

func ParseNijisanjiWikiRawTitle(title string) string {
	t := title
	t = strings.Split(t, "at")[1]
	i := strings.Index(t, " ")
	if i > 0 {
		t = t[i:]
	}
	t = strings.TrimSpace(t)
	return t
}

type ParsedNijisanjiWikiData struct {
	VideoId   string
	ChannelId string
	Handle    string
	Title     string
	Date      time.Time
}

func ParseNijisanjiWikiRawData(doc *goquery.Document) []ParsedNijisanjiWikiData {
	var parsedDatas []ParsedNijisanjiWikiData
	_ = doc.Find("h3[class*=\"date_\"]+div").Each(func(i int, s *goquery.Selection) {
		day, _ := strconv.Atoi(doc.Find("h3[class*=\"date_\"] span.day").Get(i).FirstChild.Data)
		month, _ := strconv.Atoi(doc.Find("h3[class*=\"date_\"] span.day ~ b:nth-child(3)").Get(i).FirstChild.Data)
		year, _ := strconv.Atoi(doc.Find("h3[class*=\"date_\"] span.day ~ b:nth-child(4)").Get(i).FirstChild.Data)
		s.Find("ul.list1 > li").Each(func(j int, l *goquery.Selection) {
			parsedData := ParsedNijisanjiWikiData{}
			rawText := l.Text()
			platform := GetNijisanjiLivePlatform(rawText)
			if platform != PLATFORM_YOUTUBE {
				return
			}
			splited := strings.Split(rawText, "～")
			rawDate, title := splited[0], splited[1]
			date, _ := ParseNijisanjiWikiRawDate(year, month, day, rawDate)
			l.Find("a[href^=\"https://twitter.com/\"]").Each(func(i int, a *goquery.Selection) {
				url, _ := a.Attr("href")
				twitterH := ParseTwitterUrl(url)
				if len(twitterH) < 1 {
					//fmt.Printf("twitter url parse error\n")
					return
				}
				if handle, in := NijisanjiChannelId[twitterH]; in {
					channel, err := youtube2.T.GetChannelByHandle(handle)
					if err != nil {
						return
					}
					parsedData.Handle = handle
					parsedData.ChannelId = channel.Id
				}
			})
			l.Find("a[href^=\"https://www.youtube.com/\"]").Each(func(i int, a *goquery.Selection) {
				url, _ := a.Attr("href")
				parsedData.VideoId = ParseYoutubeUrl(url)
			})
			parsedData.Title = ParseNijisanjiWikiRawTitle(title)
			parsedData.Date = date
			if len(parsedData.Handle) < 1 || len(parsedData.VideoId) < 1 || len(parsedData.ChannelId) < 1 {
				return
			}
			parsedDatas = append(parsedDatas, parsedData)
		})
	})
	return parsedDatas
}

var NijisanjiChannelId = map[string]string{
	"MahiroYukishiro": "YukishiroMahiro",
	"NaSera2434":      "NaSera",
	"ars_almal":       "ArsAlmal",
}

func RegisterNijisanjiSchedule() error {
	doc, err := RequestNijisanjiWikiRawSchedule()
	if err != nil {
		return err
	}
	schedules := ParseNijisanjiWikiRawData(doc)
	for _, s := range schedules {
		channel, err := youtube2.T.GetChannelByHandle(s.Handle)
		if err != nil {
			return err
		}
		err = AddScheduleFromVideo(AddScheduleProps{
			VideoId:    s.VideoId,
			ChannelId:  s.ChannelId,
			Handle:     s.Handle,
			Title:      s.Title,
			Thumbnail:  "",
			IsLive:     true,
			IsNowOnAir: false,
			Date:       s.Date,
		})
		if err != nil {
			return err
		}
		fmt.Printf("[info] add %v's video: %v\n", channel.Name, s.Title)
	}
	return nil
}
