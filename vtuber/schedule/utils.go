package schedule

import "strings"

func ParseYoutubeUrl(url string) string {
	u := url
	u = strings.Split(u, "=")[1]
	if strings.Contains(u, "&") {
		u = strings.Split(u, "&")[0]
	}
	return u
}

func ParseTwitterUrl(url string) string {
	u := url
	u = u[len("https://twitter.com/"):]
	u = strings.Split(u, "/")[0]
	return u
}
