package youtube

import (
	"errors"
	"google.golang.org/api/youtube/v3"
	"strconv"
	"time"
	"vdule/utils"
	"vdule/utils/cache"
)

type TubeVideo struct {
	Id                     string    `json:"id,omitempty"`
	ChannelId              string    `json:"channel_id,omitempty"`
	Title                  string    `json:"title,omitempty"`
	Description            string    `json:"description,omitempty"`
	ViewCount              uint64    `json:"view_count,omitempty"`
	LikeCount              uint64    `json:"like_count,omitempty"`
	CommentCount           uint64    `json:"comment_count,omitempty"`
	PublishedAt            string    `json:"published_at,omitempty"`
	Tags                   []string  `json:"tags,omitempty"`
	LiveBroadCastContent   string    `json:"live_broad_cast_content,omitempty"`
	IsLive                 bool      `json:"is_live,omitempty"`
	IsLiveUpcoming         bool      `json:"is_live_upcoming,omitempty"`
	IsLiveNowOnAir         bool      `json:"is_live_now_on_air,omitempty"`
	LiveActualStartTime    string    `json:"live_actual_start_time,omitempty"`
	LiveActualEndTime      string    `json:"live_actual_end_time,omitempty"`
	LiveScheduledStartTime string    `json:"live_scheduled_start_time,omitempty"`
	LiveScheduledEndTime   string    `json:"live_scheduled_end_time,omitempty"`
	LiveDate               time.Time `json:"live_date,omitempty"`
	Duration               string    `json:"duration,omitempty"`
	Thumbnail              string    `json:"thumbnail,omitempty"`
}

func GetPlayListVideosIdCache(id string, max int64) ([]string, bool) {
	cacheId := cache.GetCacheId("playlist", id, strconv.FormatInt(max, 10))
	if c, in := cache.GetCache(cacheId, &youtube.PlaylistItemListResponse{}); in {
		var videos []string
		for _, item := range c.Items {
			videos = append(videos, item.ContentDetails.VideoId)
		}
		return videos, true
	}
	return []string{}, false
}

func (t *Tube) GetPlayListVideosId(id string, max int64) ([]string, error) {
	if c, in := GetPlayListVideosIdCache(id, max); in {
		return c, nil
	}
	if max > 50 {
		return nil, errors.New("error")
	}
	call := t.Service.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(id).MaxResults(max)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	cacheId := cache.GetCacheId("playlist", id, strconv.FormatInt(max, 10))
	_ = cache.PushCacheExp(cacheId, res, time.Minute*30)
	var videos []string
	for _, item := range res.Items {
		videos = append(videos, item.ContentDetails.VideoId)
	}
	return videos, nil
}

func (t *Tube) VideoIdsToVideos(ids []string) ([]*TubeVideo, error) {
	var videos []*TubeVideo
	for _, id := range ids {
		video, err := t.GetVideo(id)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func FilterLive(videos []*TubeVideo) []*TubeVideo {
	return utils.Filter(videos, func(video *TubeVideo) bool {
		return video.IsLive
	})
}

func FilterUpcoming(videos []*TubeVideo) []*TubeVideo {
	return utils.Filter(videos, func(video *TubeVideo) bool {
		return video.IsLiveUpcoming
	})
}

func isLive(video *youtube.Video) bool {
	if video.LiveStreamingDetails != nil {
		return true
	}
	return false
}

func GetRawVideoCache(id string) (*youtube.Video, bool) {
	cacheId := cache.GetCacheId("video", id)
	if c, in := cache.GetCache(cacheId, &youtube.Video{}); in {
		return c, true
	}
	return nil, false
}

func (t *Tube) GetRawVideo(id string) (*youtube.Video, error) {
	if c, in := GetRawVideoCache(id); in {
		return c, nil
	}
	call := t.Service.Videos.List([]string{"snippet", "liveStreamingDetails", "statistics", "contentDetails"}).Id(id)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	cacheId := cache.GetCacheId("video", id)
	_ = cache.PushCache(cacheId, res.Items[0])
	return res.Items[0], nil
}

func ParseYoutubeApiDate(datetime string) time.Time {
	/*
	 * youtube data api v3の時刻はRFC3339でタイムゾーンはUTCなので
	 * 扱いやすいようにtime.Timeのjstに変換する
	 */
	t, err := time.Parse("2006-01-02T15:04:05Z", datetime)
	if err != nil {
		panic(err.Error())
	}
	t = t.Add(time.Hour * 9)
	return t
}

func VideoToTubeVideo(item *youtube.Video) TubeVideo {
	actualStartTime, actualEndTime, scheduledStartTime, scheduledEndTime, t := "", "", "", "", time.Time{}
	if item.LiveStreamingDetails != nil {
		actualStartTime = item.LiveStreamingDetails.ActualStartTime
		actualEndTime = item.LiveStreamingDetails.ActualEndTime
		scheduledStartTime = item.LiveStreamingDetails.ScheduledStartTime
		scheduledEndTime = item.LiveStreamingDetails.ScheduledEndTime
		if len(scheduledStartTime) != 0 {
			t = ParseYoutubeApiDate(scheduledStartTime)
		}
		if len(actualStartTime) != 0 {
			t = ParseYoutubeApiDate(actualStartTime)
		}
	}
	return TubeVideo{
		Id:                     item.Id,
		ChannelId:              item.Snippet.ChannelId,
		Title:                  item.Snippet.Title,
		Description:            item.Snippet.Description,
		ViewCount:              item.Statistics.ViewCount,
		LikeCount:              item.Statistics.LikeCount,
		CommentCount:           item.Statistics.CommentCount,
		PublishedAt:            item.Snippet.PublishedAt,
		Tags:                   item.Snippet.Tags,
		LiveBroadCastContent:   item.Snippet.LiveBroadcastContent,
		IsLive:                 isLive(item),
		IsLiveUpcoming:         item.Snippet.LiveBroadcastContent == "upcoming",
		IsLiveNowOnAir:         item.Snippet.LiveBroadcastContent == "live",
		LiveActualStartTime:    actualStartTime,
		LiveActualEndTime:      actualEndTime,
		LiveScheduledStartTime: scheduledStartTime,
		LiveScheduledEndTime:   scheduledEndTime,
		LiveDate:               t,
		Duration:               item.ContentDetails.Duration,
		Thumbnail:              item.Snippet.Thumbnails.High.Url,
	}
}

func (t *Tube) GetVideo(id string) (*TubeVideo, error) {
	item, err := t.GetRawVideo(id)
	if err != nil {
		return nil, err
	}
	v := VideoToTubeVideo(item)
	return &v, nil
}
