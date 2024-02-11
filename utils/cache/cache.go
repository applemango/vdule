package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"vdule/utils"
	"vdule/utils/db/redis"
)

type Id struct {
	Raw string
}

func GetCacheId(ids ...string) Id {
	return Id{Raw: utils.Reduce(ids, func(acc string, item string) string {
		return fmt.Sprintf("%v:%v", acc, item)
	}, "cache")}
}

func GetCache[T any](id Id, data *T) (*T, bool) {
	cache, err := redis.Get(id.Raw)
	if err != nil {
		return nil, false
	}
	_ = json.Unmarshal([]byte(cache), &data)
	return data, true
}

func PushCache(id Id, data any) error {
	return PushCacheExp(id, data, 0)
}

func PushCacheExp(id Id, data any, exp time.Duration) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cache := string(j)
	err = redis.PushExp(id.Raw, cache, exp)
	if err != nil {
		return err
	}
	return nil
}

func ResponseMiddleware(fn func(c *gin.Context) (int, any)) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Cache struct {
			Status int
			Json   any
		}
		url := c.Request.URL.String()
		cacheId := GetCacheId("response", url)
		if cache, in := GetCache(cacheId, &Cache{}); in {
			c.JSON(cache.Status, cache.Json)
			return
		}
		code, res := fn(c)
		_ = PushCacheExp(cacheId, Cache{
			Status: code,
			Json:   res,
		}, time.Minute)
		c.JSON(code, res)
		return
	}
}
