package youtube

import "google.golang.org/api/youtube/v3"

func (t Tube) GetChannelByHandle(handle string, parts []string) (*youtube.Channel, error) {
	call := t.Service.Channels.List(parts).ForHandle(handle)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	return res.Items[0], nil
}

func GetChannelUploadPlayList(channel youtube.Channel) string {
	return channel.ContentDetails.RelatedPlaylists.Uploads
}
