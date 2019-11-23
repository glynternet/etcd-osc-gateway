package osc

import (
	"context"
	"time"

	"github.com/glynternet/go-osc/osc"
	"github.com/pkg/errors"
)

// Dispatcher dispatches OSC Packet messages to the KeyValuePutter
// HandleError and HandleSuccess must be non-nil. Nil values for these fields will cause panics during use.
type Dispatcher struct {
	KeyValuePutter
	HandleError   func(error)
	HandleSuccess func(message osc.Message)
}

// Dispatch dispatches OSC Packet messages to the KeyValuePutter
func (d Dispatcher) Dispatch(packet osc.Packet) {
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

func (d Dispatcher) dispatchMessage(msg osc.Message) {
	k, v, err := getKeyValue(msg)
	if err != nil {
		d.HandleError(err)
		return
	}

	//TODO(glynternet): where should we get the context from?
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := d.KeyValuePutter.Put(ctx, k, v); err != nil {
		d.HandleError(errors.Wrap(err, "putting key-value pair into store"))
		return
	}
	d.HandleSuccess(msg)
}
