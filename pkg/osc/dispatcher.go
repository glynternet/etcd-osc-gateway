package osc

import (
	"context"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/pkg/errors"
)

// KeyValueDispatcher extracts key value pairs from OSC Packets passes them to the KeyValuePutter.
// HandleError and HandleSuccess must be non-nil. Nil values for these fields will cause panics during use.
type KeyValueDispatcher struct {
	KeyValuePutter
	HandleError   func(error)
	HandleSuccess func(message osc.Message, k, v string)
}

// Dispatch dispatches OSC Packet messages to the KeyValuePutter
func (d KeyValueDispatcher) Dispatch(packet osc.Packet) {
	switch p := packet.(type) {
	case *osc.Message:
		d.dispatchMessage(*p)

	case *osc.Bundle:
		for _, message := range p.Messages {
			d.dispatchMessage(*message)
		}
		for _, b := range p.Bundles {
			d.Dispatch(b)
		}
	}
}

func (d KeyValueDispatcher) dispatchMessage(msg osc.Message) {
	k, v, err := getKeyValue(msg)
	if err != nil {
		d.HandleError(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := d.KeyValuePutter.Put(ctx, k, v); err != nil {
		d.HandleError(errors.Wrap(err, "putting key-value pair into store"))
		return
	}
	d.HandleSuccess(msg, k, v)
}
