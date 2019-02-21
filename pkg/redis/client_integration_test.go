// +build integration

package redis

import (
	"github.com/klim0v/sequence-hashing/pkg/entity"
	"testing"
)

func TestClientPushPop(t *testing.T) {
	number := "111111"
	client := NewClient()
	_ = client.Push(&entity.Result{Number: number})
	recvResult, _ := client.Pop()
	_ = client.rc.Close()
	if recvResult == nil {
		t.Fatal("result is nil")
	}
	if recvResult.Number != number {
		t.Error("not equal")
	}
}
