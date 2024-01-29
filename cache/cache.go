package cache

import (
	"encoding/json"
	"fmt"
	"vdule/db/redis"
	"vdule/utils"
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
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	cache := string(j)
	err = redis.Push(id.Raw, cache)
	if err != nil {
		return err
	}
	return nil
}
