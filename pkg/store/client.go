package store

import (
	"github.com/go-redis/redis"
	"github.com/klim0v/sequence-hashing/pkg/entity"
)

const key = "list"

type Client struct {
	rc *redis.Client
}

func NewClient(rc *redis.Client) Client {
	return Client{rc}
}

func (c Client) Push(result *entity.Result) (err error) {
	err = c.rc.RPush(key, result).Err()
	return
}

func (c Client) Pop() (*entity.Result, error) {
	gotten, err := c.rc.LPop(key).Result()
	if err != nil {
		return nil, err
	}
	result := new(entity.Result)
	err = result.UnmarshalBinary([]byte(gotten))
	if err != nil {
		return nil, err
	}
	return result, nil
}
