package gojay

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeTime(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		json         string
		format       string
		err          bool
		expectedTime string
	}{
		{
			name:         "basic",
			json:         `"2018-02-18"`,
			format:       `2006-01-02`,
			err:          false,
			expectedTime: "2018-02-18",
		},
		{
			name:         "basic",
			json:         `"2017-01-02T15:04:05Z"`,
			format:       time.RFC3339,
			err:          false,
			expectedTime: "2017-01-02T15:04:05Z",
		},
		{
			name:         "basic",
			json:         `"2017-01-02T15:04:05ZINVALID"`,
			format:       time.RFC3339,
			err:          true,
			expectedTime: "",
		},
		{
			name:         "basic",
			json:         `"2017-01-02T15:04:05ZINVALID`,
			format:       time.RFC1123,
			err:          true,
			expectedTime: "",
		},
		{
			name:         "basic",
			json:         `"2017-01-02T15:04:05ZINVALID"`,
			format:       time.RFC1123,
			err:          true,
			expectedTime: "",
		},
		{
			name:         "basic",
			json:         `"2017-01-02T15:04:05ZINVALID`,
			format:       time.RFC3339,
			err:          true,
			expectedTime: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tm := time.Time{}
			dec := NewDecoder(strings.NewReader(testCase.json))
			err := dec.DecodeTime(&tm, testCase.format)
			if !testCase.err {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedTime, tm.Format(testCase.format))
				return
			}
			require.Error(t, err)
		})
	}
}

func TestDecodeAddTime(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		json         string
		format       string
		err          bool
		expectedTime string
	}{
		{
			name:         "basic",
			json:         `"2018-02-18"`,
			format:       `2006-01-02`,
			err:          false,
			expectedTime: "2018-02-18",
		},
		{
			name:         "basic",
			json:         ` "2017-01-02T15:04:05Z"`,
			format:       time.RFC3339,
			err:          false,
			expectedTime: "2017-01-02T15:04:05Z",
		},
		{
			name:         "basic",
			json:         ` "2017-01-02T15:04:05ZINVALID"`,
			format:       time.RFC3339,
			err:          true,
			expectedTime: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tm := time.Time{}
			dec := NewDecoder(strings.NewReader(testCase.json))
			err := dec.AddTime(&tm, testCase.format)
			if !testCase.err {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedTime, tm.Format(testCase.format))
				return
			}
			require.Error(t, err)
		})
	}
}

func TestDecoderTimePoolError(t *testing.T) {
	t.Parallel()

	// reset the pool to make sure it's not full
	decPool = sync.Pool{
		New: func() interface{} {
			return NewDecoder(nil)
		},
	}
	dec := NewDecoder(nil)
	dec.Release()
	defer func() {
		err := recover()
		require.Error(t, err.(error), "err shouldn't be nil")
		assert.IsType(t, InvalidUsagePooledDecoderError(""), err, "err should be of type InvalidUsagePooledDecoderError")
	}()
	_ = dec.DecodeTime(&time.Time{}, time.RFC3339)
	assert.True(t, false, "should not be called as decoder should have panicked")
}
