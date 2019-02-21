package redis

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	client2 := NewClient()
	if client != client2 {
		t.Error("not equal")
	}
}
