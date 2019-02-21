package main

import (
	"fmt"
	"testing"
)

var argTests = []struct {
	args []string
	error
}{
	{[]string{"12345678", "8"}, nil},
	{[]string{}, fmt.Errorf(argCountErrorMessage)},
	{[]string{"12345678"}, fmt.Errorf(argCountErrorMessage)},
	{[]string{"123", "8"}, fmt.Errorf(arg1ErrorMessage, "123")},
	{[]string{"aaa", "8"}, fmt.Errorf(arg1ErrorMessage, "aaa")},
	{[]string{"111111", "bb"}, fmt.Errorf(arg2ErrorMessage, "bb")},
}

func TestParse(t *testing.T) {
	for _, at := range argTests {
		t.Run(fmt.Sprint(at.args), func(t *testing.T) {
			_, _, err := Parse(at.args)
			if fmt.Sprint(at.error) != fmt.Sprint(err) {
				t.Errorf("got %q, want %q", err, at.error)
			}
		})
	}
}
