package redis

import (
	"github.com/go-redis/redis"
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"os"
	"sync"
)

var (
	once        sync.Once
	redisClient *redis.Client
)

const key = "list"

type Client struct {
	rc *redis.Client
}

func NewClient() Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDR"),
		})
	})
	return Client{redisClient}
}

func (c Client) Push(result *entity.Result) (err error) {
	err = c.rc.Publish(key, result).Err()
	return
}

func (c Client) Subscribe() *redis.PubSub {
	return c.rc.Subscribe(key)
}
