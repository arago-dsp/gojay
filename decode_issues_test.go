package gojay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeShouldNotModifyInput(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want string
	}{
		{name: "Unmarshal must not modify data", data: `"\/foo\/bar"`, want: `/foo/bar`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var s string

			data := []byte(testCase.data)
			orig := make([]byte, len(data))
			copy(orig, data)

			err := Unmarshal(data, &s)

			assert.NoError(t, err)
			assert.Equal(t, testCase.want, s)
			assert.Equal(t, orig, data)
		})
	}
}
