package gojay

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoderFloat64(t *testing.T) {
	t.Parallel()

	testCasesBasic := []struct {
		name         string
		v            float64
		expectedJSON string
	}{
		{
			name:         "basic",
			v:            float64(1),
			expectedJSON: "[1,1]",
		},
		{
			name:         "big",
			v:            float64(0),
			expectedJSON: "[0,0]",
		},
	}
	for _, testCase := range testCasesBasic {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b := &strings.Builder{}
			enc := NewEncoder(b)
			err := enc.Encode(EncodeArrayFunc(func(enc *Encoder) {
				enc.Float64(testCase.v)
				enc.AddFloat64(testCase.v)
			}))
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
	testCasesOmitEmpty := []struct {
		name         string
		v            float64
		expectedJSON string
	}{
		{
			name:         "basic",
			v:            float64(1),
			expectedJSON: "[1,1]",
		},
		{
			name:         "big",
			v:            float64(0),
			expectedJSON: "[]",
		},
	}
	for _, testCase := range testCasesOmitEmpty {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b := &strings.Builder{}
			enc := NewEncoder(b)
			err := enc.Encode(EncodeArrayFunc(func(enc *Encoder) {
				enc.Float64OmitEmpty(testCase.v)
				enc.AddFloat64OmitEmpty(testCase.v)
			}))
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
	testCasesKeyBasic := []struct {
		name         string
		v            float64
		expectedJSON string
	}{
		{
			name:         "basic",
			v:            float64(1),
			expectedJSON: `{"foo":1,"bar":1}`,
		},
		{
			name:         "big",
			v:            float64(0),
			expectedJSON: `{"foo":0,"bar":0}`,
		},
	}
	for _, testCase := range testCasesKeyBasic {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b := &strings.Builder{}
			enc := NewEncoder(b)
			err := enc.Encode(EncodeObjectFunc(func(enc *Encoder) {
				enc.Float64Key("foo", testCase.v)
				enc.AddFloat64Key("bar", testCase.v)
			}))
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
	testCasesKeyOmitEmpty := []struct {
		name         string
		v            float64
		expectedJSON string
	}{
		{
			name:         "basic",
			v:            float64(1),
			expectedJSON: `{"foo":1,"bar":1}`,
		},
		{
			name:         "big",
			v:            float64(0),
			expectedJSON: "{}",
		},
	}
	for _, testCase := range testCasesKeyOmitEmpty {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			b := &strings.Builder{}
			enc := NewEncoder(b)
			err := enc.Encode(EncodeObjectFunc(func(enc *Encoder) {
				enc.Float64KeyOmitEmpty("foo", testCase.v)
				enc.AddFloat64KeyOmitEmpty("bar", testCase.v)
			}))
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}

func TestEncoderFloat64NullEmpty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st elem",
			baseJSON:     "[",
			expectedJSON: `[null,1`,
		},
		{
			name:         "basic 2nd elem",
			baseJSON:     `["test"`,
			expectedJSON: `["test",null,1`,
		},
	}
	for _, testCase := range testCases {
		t.Run("true", func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.FloatNullEmpty(0)
			enc.AddFloatNullEmpty(1)
			_, err := enc.Write()
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}

func TestEncoderFloat64KeyNullEmpty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st elem",
			baseJSON:     "{",
			expectedJSON: `{"foo":null,"bar":1`,
		},
		{
			name:         "basic 2nd elem",
			baseJSON:     `{"test":"test"`,
			expectedJSON: `{"test":"test","foo":null,"bar":1`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.FloatKeyNullEmpty("foo", 0)
			enc.AddFloatKeyNullEmpty("bar", 1)
			_, err := enc.Write()
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}

func TestEncoderFloat32NullEmpty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st elem",
			baseJSON:     "[",
			expectedJSON: `[null,1`,
		},
		{
			name:         "basic 2nd elem",
			baseJSON:     `["test"`,
			expectedJSON: `["test",null,1`,
		},
	}
	for _, testCase := range testCases {
		t.Run("true", func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.Float32NullEmpty(0)
			enc.AddFloat32NullEmpty(1)
			_, err := enc.Write()
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}

func TestEncoderFloat32KeyNullEmpty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st elem",
			baseJSON:     "{",
			expectedJSON: `{"foo":null,"bar":1`,
		},
		{
			name:         "basic 2nd elem",
			baseJSON:     `{"test":"test"`,
			expectedJSON: `{"test":"test","foo":null,"bar":1`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.Float32KeyNullEmpty("foo", 0)
			enc.AddFloat32KeyNullEmpty("bar", 1)
			_, err := enc.Write()
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}
