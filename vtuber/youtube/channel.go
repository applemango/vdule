package youtube

import (
	"google.golang.org/api/youtube/v3"
	"strings"
	"vdule/utils/cache"
)

type TubeChannel struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Handle          string `json:"handle"`
	PublishAt       string `json:"publish_at"`
	IconImage       string `json:"icon_image"`
	UploadsPlaylist string `json:"uploads_playlist"`
	ViewCount       uint64 `json:"view_count"`
	SubscriberCount uint64 `json:"subscriber_count"`
	VideoCount      uint64 `json:"video_count"`
	Trailer         string `json:"trailer"`
	BannerImage     string `json:"banner_image"`
}

func GetRawChannelByHandleCache(handle string) (*youtube.Channel, bool) {
	cacheId := cache.GetCacheId("channel", "handle", handle)
	if c, in := cache.GetCache(cacheId, &youtube.Channel{}); in {
		return c, true
	}
	return nil, false
}

func ParseChannelHandle(handle string) string {
	p := handle
	if p[0] == '@' {
		p = p[1:]
	}
	return strings.ToLower(p)
}

func (t *Tube) GetRawChannelByHandle(handle string) (*youtube.Channel, error) {
	h := ParseChannelHandle(handle)
	if c, in := GetRawChannelByHandleCache(h); in {
		return c, nil
	}
	call := t.Service.Channels.List([]string{"snippet", "contentDetails", "brandingSettings", "statistics"}).ForHandle(h)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	cacheId := cache.GetCacheId("channel", "handle", h)
	_ = cache.PushCache(cacheId, res.Items[0])
	return res.Items[0], nil
}

func ChannelToTubeChannel(channel *youtube.Channel) TubeChannel {
	return TubeChannel{
		Id:              channel.Id,
		Name:            channel.Snippet.Title,
		Description:     channel.Snippet.Description,
		Handle:          channel.Snippet.CustomUrl[1:],
		PublishAt:       channel.Snippet.PublishedAt,
		IconImage:       channel.Snippet.Thumbnails.High.Url,
		UploadsPlaylist: channel.ContentDetails.RelatedPlaylists.Uploads,
		ViewCount:       channel.Statistics.ViewCount,
		SubscriberCount: channel.Statistics.SubscriberCount,
		VideoCount:      channel.Statistics.VideoCount,
		//Trailer:         channel.BrandingSettings.Channel.UnsubscribedTrailer,		//Trailer:         channel.BrandingSettings.Channel.UnsubscribedTrailer,
		BannerImage: channel.BrandingSettings.Image.BannerExternalUrl,
	}
}

func (t *Tube) GetChannelByHandle(handle string) (*TubeChannel, error) {
	channel, err := t.GetRawChannelByHandle(handle)
	if err != nil {
		return nil, err
	}
	c := ChannelToTubeChannel(channel)
	return &c, nil
}

func GetChannelUploadPlayList(channel *youtube.Channel) string {
	return channel.ContentDetails.RelatedPlaylists.Uploads
}

func (t *Tube) GetChannelLive(handle string) ([]*TubeVideo, error) {
	channel, err := t.GetRawChannelByHandle(handle)
	if err != nil {
		return nil, err
	}
	playList := GetChannelUploadPlayList(channel)
	videosId, err := t.GetPlayListVideosId(playList, 5)
	if err != nil {
		return nil, err
	}
	videos, err := t.VideoIdsToVideos(videosId)
	if err != nil {
		return nil, err
	}
	lives := FilterLive(videos)
	return lives, nil
}
