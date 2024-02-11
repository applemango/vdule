package redis

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var DB = ConnectRedis()

func ConnectRedis() *redis.Client {
	err := godotenv.Load(".env")
	if err != nil {
		panic("failed load config")
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("127.0.0.1:%v", os.Getenv("DB_REDIS_PORT")),
		Password: "",
		DB:       0,
	})
}

func Push(key, value string) error {
	return PushExp(key, value, 0)
}

func PushExp(key, value string, exp time.Duration) error {
	err := DB.Set(Ctx, key, value, exp).Err()
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
