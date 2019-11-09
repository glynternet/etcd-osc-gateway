package osc

import (
	"context"
	"errors"
	"testing"

	"github.com/glynternet/go-osc/osc"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher_Dispatch(t *testing.T) {
	t.Run("message should call KeyValuePutter", func(t *testing.T) {
		kvp := mockKeyValuePutter{}
		var in osc.Message
		dispatcher := Dispatcher{
			KeyValuePutter: &kvp,
			HandleSuccess:  func(msg osc.Message) { in = msg },
		}
		msg := osc.Message{Address: "woop", Arguments: []interface{}{"doop"}}
		dispatcher.Dispatch(&msg)
		assert.NotNil(t, kvp.c)
		assert.Equal(t, kvp.k, "woop")
		assert.Equal(t, kvp.v, "doop")
		assert.Equal(t, msg, in)
	})

	t.Run("bundle should call KeyValuePutter", func(t *testing.T) {
		kvp := mockKeyValuePutter{}
		dispatcher := Dispatcher{
			KeyValuePutter: &kvp,
			HandleSuccess:  func(osc.Message) {},
		}
		dispatcher.Dispatch(
			&osc.Bundle{
				Messages: []*osc.Message{
					{Address: "woop", Arguments: []interface{}{"doop"}},
					{Address: "shoop", Arguments: []interface{}{"snoop"}},
				},
				Bundles: []*osc.Bundle{{
					Messages: []*osc.Message{
						{Address: "floop", Arguments: []interface{}{"noop"}},
						{Address: "doop", Arguments: []interface{}{"boop"}},
					},
				}},
			},
		)
		assert.Equal(t, kvp.callCount, 4)
		assert.NotNil(t, kvp.c)
		assert.Equal(t, kvp.k, "doop")
		assert.Equal(t, kvp.v, "boop")
	})
}

func TestDispatcher_dispatchMessage(t *testing.T) {
	t.Run("should return key value error", func(t *testing.T) {
		var err error
		errHandler := func(inErr error) { err = inErr }
		kvp := mockKeyValuePutter{}
		dispatcher := Dispatcher{KeyValuePutter: &kvp, HandleError: errHandler}
		dispatcher.dispatchMessage(osc.Message{})
		assert.Error(t, err)
		assert.Nil(t, kvp.c)
		assert.Equal(t, "", kvp.k)
		assert.Equal(t, "", kvp.v)
	})

	t.Run("KeyValuePutter should be called", func(t *testing.T) {
		var err error
		errHandler := func(inErr error) { err = inErr }
		kvp := mockKeyValuePutter{error: errors.New("foo")}
		dispatcher := Dispatcher{KeyValuePutter: &kvp, HandleError: errHandler}
		dispatcher.dispatchMessage(osc.Message{Address: "woop", Arguments: []interface{}{"doop"}})
		assert.Error(t, err)
		assert.NotNil(t, kvp.c)
		assert.Equal(t, kvp.k, "woop")
		assert.Equal(t, kvp.v, "doop")
	})
}

type mockKeyValuePutter struct {
	callCount int
	c         context.Context
	k         string
	v         string
	error     error
}

func (m *mockKeyValuePutter) Put(c context.Context, k, v string) error {
	m.callCount++
	m.c = c
	m.k = k
	m.v = v
	return m.error
}
