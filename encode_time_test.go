package gojay

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeTime(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		tt           string
		format       string
		expectedJSON string
		err          bool
	}{
		{
			name:         "basic",
			tt:           "2018-02-01",
			format:       "2006-01-02",
			expectedJSON: `"2018-02-01"`,
			err:          false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b := strings.Builder{}
			tt, err := time.Parse(testCase.format, testCase.tt)
			require.NoError(t, err)
			enc := NewEncoder(&b)
			err = enc.EncodeTime(&tt, testCase.format)
			if !testCase.err {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedJSON, b.String())
			}
		})
	}
	t.Run("encode-time-pool-error", func(t *testing.T) {
		t.Parallel()

		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		enc.isPooled = 1
		defer func() {
			err := recover()
			require.NotNil(t, err, "err should not be nil")
			assert.IsType(t, InvalidUsagePooledEncoderError(""), err, "err should be of type InvalidUsagePooledEncoderError")
		}()
		_ = enc.EncodeTime(&time.Time{}, "")
		assert.True(t, false, "should not be called as encoder should have panicked")
	})
	t.Run("write-error", func(t *testing.T) {
		t.Parallel()

		w := TestWriterError("")
		enc := BorrowEncoder(w)
		defer enc.Release()
		err := enc.EncodeTime(&time.Time{}, "")
		require.Error(t, err)
	})
}

func TestAddTimeKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		tt           string
		format       string
		expectedJSON string
		baseJSON     string
		err          bool
	}{
		{
			name:         "basic",
			tt:           "2018-02-01",
			format:       "2006-01-02",
			baseJSON:     "{",
			expectedJSON: `{"test":"2018-02-01"`,
			err:          false,
		},
		{
			name:         "basic",
			tt:           "2018-02-01",
			format:       "2006-01-02",
			baseJSON:     `{""`,
			expectedJSON: `{"","test":"2018-02-01"`,
			err:          false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b := strings.Builder{}
			tt, err := time.Parse(testCase.format, testCase.tt)
			require.NoError(t, err)
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.AddTimeKey("test", &tt, testCase.format)
			_, err = enc.Write()
			require.NoError(t, err)
			if !testCase.err {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedJSON, b.String())
			}
		})
	}
}

func TestAddTime(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		tt           string
		format       string
		expectedJSON string
		baseJSON     string
		err          bool
	}{
		{
			name:         "basic",
			tt:           "2018-02-01",
			format:       "2006-01-02",
			baseJSON:     "[",
			expectedJSON: `["2018-02-01"`,
			err:          false,
		},
		{
			name:         "basic",
			tt:           "2018-02-01",
			format:       "2006-01-02",
			baseJSON:     "[",
			expectedJSON: `["2018-02-01"`,
			err:          false,
		},
		{
			name:         "basic",
			tt:           "2018-02-01",
			format:       "2006-01-02",
			baseJSON:     `[""`,
			expectedJSON: `["","2018-02-01"`,
			err:          false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b := strings.Builder{}
			tt, err := time.Parse(testCase.format, testCase.tt)
			require.NoError(t, err)
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.AddTime(&tt, testCase.format)
			_, err = enc.Write()
			require.NoError(t, err)
			if !testCase.err {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedJSON, b.String())
			}
		})
	}
}
