package youtube

import (
	"errors"
	"google.golang.org/api/youtube/v3"
)

func (t *Tube) GetPlayListVideosId(id string, max int64) ([]string, error) {
	if max > 50 {
		return nil, errors.New("error")
	}
	call := t.Service.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(id).MaxResults(max)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
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
	var res []*TubeVideo
	for _, video := range videos {
		if video.isLive {
			res = append(res, video)
		}
	}
	return res
}

type TubeVideo struct {
	Id                     string
	ChannelId              string
	Title                  string
	Description            string
	ViewCount              uint64
	LikeCount              uint64
	CommentCount           uint64
	PublishedAt            string
	Tags                   []string
	LiveBroadCastContent   string
	isLive                 bool
	isLiveUpcoming         bool
	LiveActualStartTime    string
	LiveActualEndTime      string
	LiveScheduledStartTime string
	LiveScheduledEndTime   string

	Duration string
}

func isLive(video *youtube.Video) bool {
	if video.LiveStreamingDetails != nil {
		return true
	}
	return false
}

func (t *Tube) GetRawVideo(id string) (*youtube.Video, error) {
	call := t.Service.Videos.List([]string{"snippet", "liveStreamingDetails", "statistics", "contentDetails"}).Id(id)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	return res.Items[0], nil
}

func (t *Tube) GetVideo(id string) (*TubeVideo, error) {
	item, err := t.GetRawVideo(id)
	if err != nil {
		return nil, err
	}

	actualStartTime, actualEndTime, scheduledStartTime, scheduledEndTime := "", "", "", ""
	if item.LiveStreamingDetails != nil {
		actualStartTime = item.LiveStreamingDetails.ActualStartTime
		actualEndTime = item.LiveStreamingDetails.ActualEndTime
		scheduledStartTime = item.LiveStreamingDetails.ScheduledStartTime
		scheduledEndTime = item.LiveStreamingDetails.ScheduledEndTime
	}

	return &TubeVideo{
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
		isLive:                 isLive(item),
		isLiveUpcoming:         item.Snippet.LiveBroadcastContent == "live",
		LiveActualStartTime:    actualStartTime,
		LiveActualEndTime:      actualEndTime,
		LiveScheduledStartTime: scheduledStartTime,
		LiveScheduledEndTime:   scheduledEndTime,
		Duration:               item.ContentDetails.Duration,
	}, nil
}
