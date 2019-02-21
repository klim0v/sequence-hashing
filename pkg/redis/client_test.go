package redis

import (
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	client2 := NewClient()
	if client != client2 {
		t.Error("not equal")
	}
}

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
