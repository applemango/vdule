package youtube

import "errors"

func (t Tube) GetPlayListVideos(id string, max int64) ([]string, error) {
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

func (t Tube) GetVideo(id string) (error, error) {
	call := t.Service.Videos.List([]string{"snippet"}).Id(id)
	res, err := call.Do()
	if err != nil {
		return nil, err
	}
	_ = res.Items[0]
	return nil, nil
}
