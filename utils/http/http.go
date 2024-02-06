package http

import (
	"github.com/wenerme/go-req"
	"io"
	"net/http"
	"time"
	"vdule/utils/cache"
)

/*func PushCache(url string, body string) error {
	err := redis.Push(fmt.Sprintf("cache:%s", url), body)
	return err
}

func GetCache(url string) (string, error) {
	body, err := redis.Get(fmt.Sprintf("cache:%s", url))
	return body, err
}*/

func Request(url string, method string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "ja,en-US;q=0.9,en;q=0.8,ja-JP;q=0.7")
	req.Header.Add("sec-ch-ua", `"Not/A)Brand";v="99", "Google Chrome";v="115", "Chromium";v="115"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"macOS"`)
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")

	//req.Header.Set("authority", "schedule.hololive.tv")
	//req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//req.Header.Set("accept-language", "ja,en;q=0.9,en-GB;q=0.8,en-US;q=0.7")
	//req.Header.Set("cache-control", "max-age=0")
	//req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Microsoft Edge";v="121", "Chromium";v="121"`)
	//req.Header.Set("sec-ch-ua-mobile", "?0")
	//req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	//req.Header.Set("sec-fetch-dest", "document")
	//req.Header.Set("sec-fetch-mode", "navigate")
	//req.Header.Set("sec-fetch-site", "none")
	//req.Header.Set("sec-fetch-user", "?1")
	//req.Header.Set("upgrade-insecure-requests", "1")
	//req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0")
	client := new(http.Client)
	resp, err := client.Do(req)
	return resp, err
}

/*func HttpGetStruct[T any](url string, obj *T) error {
	res, err := Request(url, "get", nil)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	result := string(bodyBytes)
	fmt.Printf("%v", result)
	err = json.Unmarshal([]byte(result), &obj)
	return err
}*/

func HttpGetStruct[T any](url string, obj *T) error {
	client := req.Request{
		BaseURL: url,
		Options: []interface{}{req.JSONEncode, req.JSONDecode},
	}
	//client = client.WithHook(req.DebugHook(&req.DebugOptions{}))
	err := client.With(req.Request{
		Method: http.MethodGet,
	}).Fetch(&obj)
	return err
}

func HttpGetStructCache(url string, obj *any) error {
	return HttpGetStructCacheExp(url, obj, 0)
}

func HttpGetStructCacheExp[T any](url string, obj *T, exp time.Duration) error {
	cacheId := cache.GetCacheId("http", url)
	if _, in := cache.GetCache(cacheId, &obj); in {
		return nil
	}
	err := HttpGetStruct(url, obj)
	if err != nil {
		return err
	}
	_ = cache.PushCacheExp(cacheId, obj, exp)
	return nil

}

/*type RequestOption struct {
	Data   any
	Url    string
	Method string
	Body   io.Reader
	Cache  bool
}

func RequestJSON(option RequestOption) error {
	cacheId := cache.GetCacheId("http", option.Url)
	if option.Cache {
		err := cache.GetCache(cacheId, &option.Data)
		if err == nil {
			json.Unmarshal([]byte(cache), &option.Data)
			return nil
		}
	}
	res, err := Request(option.Url, option.Method, option.Body)
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	result := string(bodyBytes)
	PushCache(option.Url, result)
	json.Unmarshal([]byte(result), &option.Data)
	return nil
}*/
