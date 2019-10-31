package cmd_test

import (
	"testing"

	"github.com/glynternet/pkg/cmd"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	bytes []byte
	error
}

func (w *mockWriter) Write(bs []byte) (int, error) {
	w.bytes = bs
	return 0, w.error
}

func TestNew(t *testing.T) {
	for _, test := range []struct {
		name    string
		version string
	}{
		{
			name: "zero values",
		},
		{
			name:    "filled version",
			version: "vWoop.Woop.Woop",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var w mockWriter
			cmd := cmd.NewVersionCmd(test.version, &w)
			err := cmd.RunE(nil, nil)
			assert.NoError(t, err)
			assert.Equal(t, []byte(test.version), w.bytes)
		})
	}

	t.Run("error", func(t *testing.T) {
		mockErr := errors.New("write error")
		w := mockWriter{error: mockErr}
		cmd := cmd.NewVersionCmd("anything", &w)
		err := cmd.RunE(nil, nil)
		assert.Equal(t, mockErr, err)
	})
}
