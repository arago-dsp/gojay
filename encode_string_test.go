package gojay

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoderStringEncodeAPI(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString("hello world")
		require.NoError(t, err)
		assert.Equal(
			t,
			`"hello world"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("utf8", func(t *testing.T) {
		t.Parallel()

		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString("漢字𩸽")
		require.NoError(t, err)
		assert.Equal(
			t,
			`"漢字𩸽"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("utf8-multibyte", func(t *testing.T) {
		t.Parallel()

		str := "テュールスト マーティン ヤコブ 😁"
		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString(str)
		require.NoError(t, err)
		assert.Equal(
			t,
			`"テュールスト マーティン ヤコブ 😁"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("escaped-sequence1", func(t *testing.T) {
		t.Parallel()

		str := `テュールスト マ\ーテ
ィン ヤコブ 😁`
		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString(str)
		require.NoError(t, err)
		assert.Equal(
			t,
			`"テュールスト マ\\ーテ\nィン ヤコブ 😁"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("escaped-sequence2", func(t *testing.T) {
		t.Parallel()

		str := `テュールスト マ\ーテ
ィン ヤコブ 😁	`
		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString(str)
		require.NoError(t, err)
		assert.Equal(
			t,
			`"テュールスト マ\\ーテ\nィン ヤコブ 😁\t"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("escaped-sequence3", func(t *testing.T) {
		t.Parallel()

		str := "hello \r world 𝄞"
		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString(str)
		require.NoError(t, err)
		assert.Equal(
			t,
			`"hello \r world 𝄞"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("escaped-sequence3", func(t *testing.T) {
		t.Parallel()

		str := "hello \b world 𝄞"
		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString(str)
		require.NoError(t, err)
		assert.Equal(
			t,
			`"hello \b world 𝄞"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("escaped-control-char", func(t *testing.T) {
		t.Parallel()

		str := "\u001b"
		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString(str)
		require.NoError(t, err)
		assert.Equal(
			t,
			`"\u001b"`,
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
	t.Run("escaped-sequence3", func(t *testing.T) {
		t.Parallel()

		str := "hello \f world 𝄞"
		builder := &strings.Builder{}
		enc := NewEncoder(builder)
		err := enc.EncodeString(str)
		require.NoError(t, err)
		assert.Equal(
			t,
			"\"hello \\f world 𝄞\"",
			builder.String(),
			"Result of marshalling is different as the one expected")
	})
}

func TestEncoderStringEncodeAPIErrors(t *testing.T) {
	t.Parallel()

	t.Run("pool-error", func(t *testing.T) {
		t.Parallel()

		v := ""
		enc := BorrowEncoder(nil)
		enc.isPooled = 1
		defer func() {
			err := recover()
			require.NotNil(t, err, "err should not be nil")
			assert.IsType(t, InvalidUsagePooledEncoderError(""), err, "err should be of type InvalidUsagePooledEncoderError")
			assert.Equal(t, "Invalid usage of pooled encoder", err.(InvalidUsagePooledEncoderError).Error(), "err should be of type InvalidUsagePooledDecoderError")
		}()
		_ = enc.EncodeString(v)
		assert.True(t, false, "should not be called as it should have panicked")
	})
	t.Run("write-error", func(t *testing.T) {
		t.Parallel()

		v := "test"
		w := TestWriterError("")
		enc := BorrowEncoder(w)
		defer enc.Release()
		err := enc.EncodeString(v)
		require.Error(t, err)
	})
}

func TestEncoderStringMarshalAPI(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		r, err := Marshal("string")
		require.NoError(t, err)
		assert.Equal(
			t,
			`"string"`,
			string(r),
			"Result of marshalling is different as the one expected")
	})
	t.Run("utf8", func(t *testing.T) {
		t.Parallel()

		r, err := Marshal("漢字")
		require.NoError(t, err)
		assert.Equal(
			t,
			`"漢字"`,
			string(r),
			"Result of marshalling is different as the one expected")
	})
}

func TestEncoderStringNullEmpty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st elem",
			baseJSON:     "[",
			expectedJSON: `[null,"true"`,
		},
		{
			name:         "basic 2nd elem",
			baseJSON:     `["test"`,
			expectedJSON: `["test",null,"true"`,
		},
		{
			name:         "basic 2nd elem",
			baseJSON:     `["test"`,
			expectedJSON: `["test",null,"true"`,
		},
	}
	for _, testCase := range testCases {
		t.Run("true", func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.StringNullEmpty("")
			enc.AddStringNullEmpty("true")
			_, err := enc.Write()
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}

func TestEncoderStringNullEmpty2(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st elem",
			baseJSON:     "[",
			expectedJSON: `["test"`,
		},
	}
	for _, testCase := range testCases {
		t.Run("true", func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.StringNullEmpty("test")
			_, err := enc.Write()
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}

func TestEncoderStringNullKeyEmpty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		baseJSON     string
		expectedJSON string
	}{
		{
			name:         "basic 1st elem",
			baseJSON:     "{",
			expectedJSON: `{"foo":null,"bar":"true"`,
		},
		{
			name:         "basic 2nd elem",
			baseJSON:     `{"test":"test"`,
			expectedJSON: `{"test":"test","foo":null,"bar":"true"`,
		},
	}
	for _, testCase := range testCases {
		t.Run("true", func(t *testing.T) {
			t.Parallel()

			var b strings.Builder
			enc := NewEncoder(&b)
			enc.writeString(testCase.baseJSON)
			enc.StringKeyNullEmpty("foo", "")
			enc.AddStringKeyNullEmpty("bar", "true")
			_, err := enc.Write()
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedJSON, b.String())
		})
	}
}
