package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var DB = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6381",
	Password: "",
	DB:       0,
})

func Push(key string, value string) error {
	err := DB.Set(Ctx, key, value, 0).Err()
	return err
}

func Get(key string) (string, error) {
	v, err := DB.Get(Ctx, key).Result()
	return v, err
}

func SearchKV(key string, max int64) ([]string, []string, error) {
	var keys []string
	var values []string
	var i int64
	iter := DB.Scan(Ctx, 0, key, max).Iterator()
	for iter.Next(Ctx) {
		if i >= max {
			return keys, values, nil
		}
		i++
		val := iter.Val()
		v, err := Get(val)
		if err != nil {
			continue
		}
		keys = append(keys, val)
		values = append(values, v)
	}
	return keys, values, nil
}
