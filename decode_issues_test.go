package gojay

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeShouldNotModifyInput(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		data string
		want string
	}{
		{name: "Unmarshal must not modify data", data: `"\/foo\/bar"`, want: `/foo/bar`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var s string

			data := []byte(testCase.data)
			orig := make([]byte, len(data))
			copy(orig, data)

			err := Unmarshal(data, &s)

			require.NoError(t, err)
			assert.Equal(t, testCase.want, s)
			assert.Equal(t, orig, data)
		})
	}
}
