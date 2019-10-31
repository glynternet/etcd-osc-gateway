package osc

import (
	"context"
	"fmt"

	"github.com/glynternet/go-osc/osc"
)

// KeyValuePutter puts a key and value string
type KeyValuePutter interface {
	Put(context.Context, string, string) error
}

func getKeyValue(msg osc.Message) (string, string, error) {
	lenArgs := len(msg.Arguments)
	if lenArgs != 1 {
		return "", "", fmt.Errorf("expected 1 arg but got %d", lenArgs)
	}
	val, ok := msg.Arguments[0].(string)
	if !ok {
		return "", "", fmt.Errorf("expected string arg but got %T: %+v", msg.Arguments[0], msg)
	}
	return msg.Address, val, nil
}
