package osc

import (
	"context"
	"errors"
	"testing"

	"github.com/hypebeast/go-osc/osc"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher_Dispatch(t *testing.T) {
	t.Run("message should call KeyValuePutter", func(t *testing.T) {
		kvp := mockKeyValuePutter{}
		var successHandleMsg osc.Message
		var successHandleK, successHandleV string
		dispatcher := KeyValueDispatcher{
			KeyValuePutter: &kvp,
			HandleSuccess: func(msg osc.Message, k, v string) {
				successHandleMsg = msg
				successHandleK = k
				successHandleV = v
			},
		}
		msg := osc.Message{Address: "woop", Arguments: []interface{}{"doop"}}
		dispatcher.Dispatch(&msg)
		assert.NotNil(t, kvp.c)
		assert.Equal(t, "woop", kvp.k)
		assert.Equal(t, "doop", kvp.v)
		assert.Equal(t, msg, successHandleMsg)
		assert.Equal(t, "woop", successHandleK)
		assert.Equal(t, "doop", successHandleV)
	})

	t.Run("bundle should call KeyValuePutter", func(t *testing.T) {
		kvp := mockKeyValuePutter{}
		dispatcher := KeyValueDispatcher{
			KeyValuePutter: &kvp,
			HandleSuccess:  func(osc.Message, string, string) {},
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
		dispatcher := KeyValueDispatcher{KeyValuePutter: &kvp, HandleError: errHandler}
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
		dispatcher := KeyValueDispatcher{KeyValuePutter: &kvp, HandleError: errHandler}
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
