package youtube

import "google.golang.org/api/youtube/v3"

func (t *Tube) GetChannelByHandle(handle string) (*youtube.Channel, error) {
	call := t.Service.Channels.List([]string{"snippet", "contentDetails"}).ForHandle(handle)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
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
