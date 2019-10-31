package osc

import (
	"context"

	"github.com/glynternet/go-osc/osc"
	"github.com/pkg/errors"
)

// Dispatcher dispatches OSC Packet messages to the KeyValuePutter
type Dispatcher struct {
	KeyValuePutter
	HandleError func(error)
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
	if err := d.KeyValuePutter.Put(context.TODO(), k, v); err != nil {
		d.HandleError(errors.Wrap(err, "putting key-value pair into store"))
	}
}
