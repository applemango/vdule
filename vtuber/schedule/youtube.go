package schedule

import (
	"fmt"
	"vdule/vtuber"
	"vdule/vtuber/youtube"
)

func RegisterYoutuberSchedule() error {
	vs, err := vtuber.GetVtubersCore(`SELECT id, handle, name, description, organization_id, is_crawl FROM vtuber WHERE is_crawl = true`)
	if err != nil {
		return err
	}
	for _, v := range vs {
		if channel, in := youtube.GetRawChannelByHandleCache(v.Handle); in {
			videos := youtube.GetChannelUploadVideos(channel)
			fmt.Printf("[info] scan %v\n", channel.Snippet.Title)
			for videos.Next() {
				v, _ := videos.Get()
				video := youtube.VideoToTubeVideo(v)
				if !(video.IsLiveUpcoming || video.IsLive) {
					break
				}
				fmt.Printf("[info] add %v's video: find %v\n", channel.Snippet.Title, v.Id)
				err = AddScheduleFromRawVideo(*channel, *v)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
