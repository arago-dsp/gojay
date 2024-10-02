package gojay

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeNumberExra(t *testing.T) {
	t.Parallel()

	t.Run("skip-number-err", func(t *testing.T) {
		t.Parallel()

		dec := NewDecoder(strings.NewReader("123456afzfz343"))
		_, err := dec.skipNumber()
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
	t.Run("get-exponent-err", func(t *testing.T) {
		t.Parallel()

		v := 0
		dec := NewDecoder(strings.NewReader("1.2Ea"))
		err := dec.Decode(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}
