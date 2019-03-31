// +build integration

package store

import (
	"github.com/go-redis/redis"
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"testing"
)

var storeClient = NewClient(redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_ADDR"),
}))

func TestClientPushPop(t *testing.T) {
	number := "111111"
	_ = storeClient.Push(&entity.Result{Number: number})
	recvResult, _ := client.Pop()
	if recvResult == nil {
		t.Fatal("result is nil")
	}
	if recvResult.Number != number {
		t.Error("not equal")
	}
}

func TestClientPopUnmarshal(t *testing.T) {
	storeClient.rc.RPush(key, "test")
	_, err := client.Pop()
	if err == redis.Nil {
		t.Error(err)
	}
}

func TestClientPopEmpty(t *testing.T) {
	storeClient.rc.Del(key)
	_, err := client.Pop()
	if err != nil && err != redis.Nil {
		t.Error(err)
	}
}
