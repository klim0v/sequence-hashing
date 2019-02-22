// +build integration

package redis

import (
	"github.com/go-redis/redis"
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"testing"
)

func TestClientPushPop(t *testing.T) {
	number := "111111"
	client := NewClient()
	_ = client.Push(&entity.Result{Number: number})
	recvResult, _ := client.Pop()
	if recvResult == nil {
		t.Fatal("result is nil")
	}
	if recvResult.Number != number {
		t.Error("not equal")
	}
}

func TestClientPopUnmarshal(t *testing.T) {
	client := NewClient()
	client.rc.RPush(key, "test")
	_, err := client.Pop()
	if err == redis.Nil {
		t.Error(err)
	}
}

func TestClientPopEmpty(t *testing.T) {
	client := NewClient()
	client.rc.Del(key)
	_, err := client.Pop()
	if err != nil && err != redis.Nil {
		t.Error(err)
	}
}
