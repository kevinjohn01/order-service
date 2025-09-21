package cache

import (
    "context"
    "time"
    "github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() {
    Rdb = redis.NewClient(&redis.Options{
        Addr: "redis:6379",
        DB:   0,
    })
}

func SetCache(key string, value string, ttl time.Duration) error {
    return Rdb.Set(Ctx, key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
    return Rdb.Get(Ctx, key).Result()
}
