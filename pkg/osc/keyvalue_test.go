package osc

import (
	"testing"

	"github.com/glynternet/go-osc/osc"
	"github.com/stretchr/testify/assert"
)

func Test_getKeyValue(t *testing.T) {
	for _, tc := range []struct {
		name      string
		msg       osc.Message
		k, v      string
		expectErr bool
	}{
		{
			name:      "wrong arg count",
			msg:       osc.Message{},
			expectErr: true,
		},
		{
			name:      "wrong arg type",
			msg:       osc.Message{Arguments: []interface{}{0}},
			expectErr: true,
		},
		{
			name: "valid message",
			msg:  osc.Message{Address: "woop", Arguments: []interface{}{"doop"}},
			k:    "woop",
			v:    "doop",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k, v, err := getKeyValue(tc.msg)
			assert.Equal(t, tc.k, k)
			assert.Equal(t, tc.v, v)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
