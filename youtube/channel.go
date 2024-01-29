package youtube

import (
	"google.golang.org/api/youtube/v3"
	"vdule/cache"
)

func GetChannelByHandleCache(handle string) (*youtube.Channel, bool) {
	cacheId := cache.GetCacheId("channel", "handle", handle)
	if c, in := cache.GetCache(cacheId, &youtube.Channel{}); in {
		return c, true
	}
	return nil, false
}

func (t *Tube) GetChannelByHandle(handle string) (*youtube.Channel, error) {
	if c, in := GetChannelByHandleCache(handle); in {
		return c, nil
	}
	call := t.Service.Channels.List([]string{"snippet", "contentDetails", "brandingSettings"}).ForHandle(handle)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	cacheId := cache.GetCacheId("channel", "handle", handle)
	_ = cache.PushCache(cacheId, res.Items[0])
	return res.Items[0], nil
}

func GetChannelUploadPlayList(channel *youtube.Channel) string {
	return channel.ContentDetails.RelatedPlaylists.Uploads
}

func (t *Tube) GetChannelLive(handle string) ([]*TubeVideo, error) {
	channel, err := t.GetChannelByHandle(handle)
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
