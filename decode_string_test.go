package gojay

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoderString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult string
		err            bool
		errType        any
	}{
		{
			name:           "basic-string",
			json:           `"string"`,
			expectedResult: "string",
			err:            false,
		},
		{
			name:           "string-solidus",
			json:           `"\/"`,
			expectedResult: "/",
			err:            false,
		},
		{
			name:           "basic-string",
			json:           ``,
			expectedResult: "",
			err:            false,
		},
		{
			name:           "basic-string",
			json:           `""`,
			expectedResult: "",
			err:            false,
		},
		{
			name:           "basic-string2",
			json:           `"hello world!"`,
			expectedResult: "hello world!",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\n"`,
			expectedResult: "\n",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\n"`,
			expectedResult: `\n`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\t"`,
			expectedResult: "\t",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\t"`,
			expectedResult: `\t`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\b"`,
			expectedResult: "\b",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\b"`,
			expectedResult: `\b`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\f"`,
			expectedResult: "\f",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\f"`,
			expectedResult: `\f`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\r"`,
			expectedResult: "\r",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "escape-control-char-solidus",
			json:           `"\/"`,
			expectedResult: "/",
			err:            false,
		},
		{
			name:           "escape-control-char-solidus",
			json:           `"/"`,
			expectedResult: "/",
			err:            false,
		},
		{
			name:           "escape-control-char-solidus-escape-char",
			json:           `"\\/"`,
			expectedResult: `\/`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\r"`,
			expectedResult: `\r`,
			err:            false,
		},
		{
			name:           "utf8",
			json:           `"𠜎 𠜱 𠝹 𠱓 𠱸 𠲖 𠳏 𠳕 𠴕 𠵼 𠵿"`,
			expectedResult: "𠜎 𠜱 𠝹 𠱓 𠱸 𠲖 𠳏 𠳕 𠴕 𠵼 𠵿",
			err:            false,
		},
		{
			name:           "utf8-code-point",
			json:           `"\u06fc"`,
			expectedResult: `ۼ`,
			err:            false,
		},
		{
			name:           "utf8-code-point-escaped",
			json:           `"\\u2070"`,
			expectedResult: `\u2070`,
			err:            false,
		},
		{
			name:           "utf8-code-point-err",
			json:           `"\u2Z70"`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\uDD1E"`,
			expectedResult: `𝄞`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\\"`,
			expectedResult: `�\`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\uD834"`,
			expectedResult: "�\x00\x00\x00",
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834"`,
			expectedResult: `�`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-err",
			json:           `"\uD834\`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-err2",
			json:           `"\uD834\uDZ1E`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-err3",
			json:           `"\uD834`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\t"`,
			expectedResult: "�\t",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\n"`,
			expectedResult: "�\n",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\f"`,
			expectedResult: "�\f",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\b"`,
			expectedResult: "�\b",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\r"`,
			expectedResult: "�\r",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\h"`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "null",
			json:           `null`,
			expectedResult: "",
		},
		{
			name:           "null-err",
			json:           `nall`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "escape quote err",
			json:           `"test string \" escaped"`,
			expectedResult: `test string " escaped`,
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \t escaped"`,
			expectedResult: "test string \t escaped",
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \r escaped"`,
			expectedResult: "test string \r escaped",
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \b escaped"`,
			expectedResult: "test string \b escaped",
			err:            false,
		},
		{
			name:           "escape quote err",
			json:           `"test string \n escaped"`,
			expectedResult: "test string \n escaped",
			err:            false,
		},
		{
			name:           "escape quote err",
			json:           `"test string \\\" escaped`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "escape quote err",
			json:           `"test string \\\l escaped"`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-json",
			json:           `invalid`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "string-complex",
			json:           `  "string with spaces and \"escape\"d \"quotes\" and escaped line returns \n and escaped \\\\ escaped char"`,
			expectedResult: "string with spaces and \"escape\"d \"quotes\" and escaped line returns \n and escaped \\\\ escaped char",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			str := ""
			dec := NewDecoder(strings.NewReader(testCase.json))
			err := dec.Decode(&str)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should of the given type")
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(
				t,
				testCase.expectedResult,
				str,
				"'%s' should be equal to expectedResult",
				str,
			)
		})
	}
}

func TestDecoderStringNull(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult string
		err            bool
		errType        any
		resultIsNil    bool
	}{
		{
			name:           "basic-string",
			json:           `"string"`,
			expectedResult: "string",
			err:            false,
		},
		{
			name:           "string-solidus",
			json:           `"\/"`,
			expectedResult: "/",
			err:            false,
		},
		{
			name:           "basic-string",
			json:           ``,
			expectedResult: "",
			err:            false,
			resultIsNil:    true,
		},
		{
			name:           "basic-string",
			json:           `""`,
			expectedResult: "",
			err:            false,
		},
		{
			name:           "basic-string2",
			json:           `"hello world!"`,
			expectedResult: "hello world!",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\n"`,
			expectedResult: "\n",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\n"`,
			expectedResult: `\n`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\t"`,
			expectedResult: "\t",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\t"`,
			expectedResult: `\t`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\b"`,
			expectedResult: "\b",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\b"`,
			expectedResult: `\b`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\f"`,
			expectedResult: "\f",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\f"`,
			expectedResult: `\f`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\r"`,
			expectedResult: "\r",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "escape-control-char-solidus",
			json:           `"\/"`,
			expectedResult: "/",
			err:            false,
		},
		{
			name:           "escape-control-char-solidus",
			json:           `"/"`,
			expectedResult: "/",
			err:            false,
		},
		{
			name:           "escape-control-char-solidus-escape-char",
			json:           `"\\/"`,
			expectedResult: `\/`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\r"`,
			expectedResult: `\r`,
			err:            false,
		},
		{
			name:           "utf8",
			json:           `"𠜎 𠜱 𠝹 𠱓 𠱸 𠲖 𠳏 𠳕 𠴕 𠵼 𠵿"`,
			expectedResult: "𠜎 𠜱 𠝹 𠱓 𠱸 𠲖 𠳏 𠳕 𠴕 𠵼 𠵿",
			err:            false,
		},
		{
			name:           "utf8-code-point",
			json:           `"\u06fc"`,
			expectedResult: `ۼ`,
			err:            false,
		},
		{
			name:           "utf8-code-point-escaped",
			json:           `"\\u2070"`,
			expectedResult: `\u2070`,
			err:            false,
		},
		{
			name:           "utf8-code-point-err",
			json:           `"\u2Z70"`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\uDD1E"`,
			expectedResult: `𝄞`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\\"`,
			expectedResult: `�\`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\uD834"`,
			expectedResult: "�\x00\x00\x00",
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834"`,
			expectedResult: `�`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-err",
			json:           `"\uD834\`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-err2",
			json:           `"\uD834\uDZ1E`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-err3",
			json:           `"\uD834`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\t"`,
			expectedResult: "�\t",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\n"`,
			expectedResult: "�\n",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\f"`,
			expectedResult: "�\f",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\b"`,
			expectedResult: "�\b",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\r"`,
			expectedResult: "�\r",
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\h"`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "null",
			json:           `null`,
			expectedResult: "",
			resultIsNil:    true,
		},
		{
			name:           "null-err",
			json:           `nall`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "escape quote err",
			json:           `"test string \" escaped"`,
			expectedResult: `test string " escaped`,
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \t escaped"`,
			expectedResult: "test string \t escaped",
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \r escaped"`,
			expectedResult: "test string \r escaped",
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \b escaped"`,
			expectedResult: "test string \b escaped",
			err:            false,
		},
		{
			name:           "escape quote err",
			json:           `"test string \n escaped"`,
			expectedResult: "test string \n escaped",
			err:            false,
		},
		{
			name:           "escape quote err",
			json:           `"test string \\\" escaped`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "escape quote err",
			json:           `"test string \\\l escaped"`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-json",
			json:           `invalid`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "string-complex",
			json:           `  "string with spaces and \"escape\"d \"quotes\" and escaped line returns \n and escaped \\\\ escaped char"`,
			expectedResult: "string with spaces and \"escape\"d \"quotes\" and escaped line returns \n and escaped \\\\ escaped char",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			str := (*string)(nil)
			err := Unmarshal([]byte(testCase.json), &str)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should of the given type")
				}
				return
			}
			require.NoError(t, err)
			if testCase.resultIsNil {
				assert.Nil(t, str)
			} else {
				assert.Equal(
					t,
					testCase.expectedResult,
					*str,
					"v must be equal to %s",
					testCase.expectedResult,
				)
			}
		})
	}
	t.Run("decoder-api-invalid-json2", func(t *testing.T) {
		t.Parallel()

		v := new(string)
		dec := NewDecoder(strings.NewReader(`a`))
		err := dec.StringNull(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderStringNoEscape(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult string
		err            bool
		errType        any
	}{
		{
			name:           "basic-string",
			json:           `"string"`,
			expectedResult: "string",
			err:            false,
		},
		{
			name:           "string-solidus",
			json:           `"\/"`,
			expectedResult: "\\/",
			err:            false,
		},
		{
			name:           "basic-string",
			json:           ``,
			expectedResult: "",
			err:            false,
		},
		{
			name:           "basic-string",
			json:           `""`,
			expectedResult: "",
			err:            false,
		},
		{
			name:           "basic-string2",
			json:           `"hello world!"`,
			expectedResult: "hello world!",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\n"`,
			expectedResult: "\\n",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\n"`,
			expectedResult: `\\n`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\t"`,
			expectedResult: "\\t",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\t"`,
			expectedResult: `\\t`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\b"`,
			expectedResult: "\\b",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\b"`,
			expectedResult: `\\b`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\f"`,
			expectedResult: "\\f",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\f"`,
			expectedResult: `\\f`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\r"`,
			expectedResult: "\\r",
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "escape-control-char-solidus",
			json:           `"\/"`,
			expectedResult: "\\/",
			err:            false,
		},
		{
			name:           "escape-control-char-solidus",
			json:           `"/"`,
			expectedResult: "/",
			err:            false,
		},
		{
			name:           "escape-control-char-solidus-escape-char",
			json:           `"\\/"`,
			expectedResult: `\\/`,
			err:            false,
		},
		{
			name:           "escape-control-char",
			json:           `"\\r"`,
			expectedResult: `\\r`,
			err:            false,
		},
		{
			name:           "utf8",
			json:           `"𠜎 𠜱 𠝹 𠱓 𠱸 𠲖 𠳏 𠳕 𠴕 𠵼 𠵿"`,
			expectedResult: "𠜎 𠜱 𠝹 𠱓 𠱸 𠲖 𠳏 𠳕 𠴕 𠵼 𠵿",
			err:            false,
		},
		{
			name:           "utf8-code-point",
			json:           `"\u06fc"`,
			expectedResult: `\u06fc`,
			err:            false,
		},
		{
			name:           "utf8-code-point-escaped",
			json:           `"\\u2070"`,
			expectedResult: `\\u2070`,
			err:            false,
		},
		{
			name:           "utf8-code-point-err",
			json:           `"\u2Z70"`,
			expectedResult: `\u2Z70`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\uDD1E"`,
			expectedResult: `\uD834\uDD1E`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\\"`,
			expectedResult: `\uD834\\`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834\uD834"`,
			expectedResult: `\uD834\uD834`,
			err:            false,
		},
		{
			name:           "utf16-surrogate",
			json:           `"\uD834"`,
			expectedResult: `\uD834`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-err",
			json:           `"\uD834\`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-err2",
			json:           `"\uD834\uDZ1E`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-err3",
			json:           `"\uD834`,
			expectedResult: ``,
			err:            true,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\t"`,
			expectedResult: `\uD834\t`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\n"`,
			expectedResult: `\uD834\n`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\f"`,
			expectedResult: `\uD834\f`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\b"`,
			expectedResult: `\uD834\b`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\r"`,
			expectedResult: `\uD834\r`,
			err:            false,
		},
		{
			name:           "utf16-surrogate-followed-by-control-char",
			json:           `"\uD834\h"`,
			expectedResult: `\uD834\h`,
			err:            false,
		},
		{
			name:           "null",
			json:           `null`,
			expectedResult: "",
		},
		{
			name:           "null-err",
			json:           `nall`,
			expectedResult: "",
			err:            true,
		},
		{
			name:           "escape quote err",
			json:           `"test string \" escaped"`,
			expectedResult: `test string \`,
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \t escaped"`,
			expectedResult: "test string \\t escaped",
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \r escaped"`,
			expectedResult: "test string \\r escaped",
			err:            false,
		},
		{
			name:           "escape quote err2",
			json:           `"test string \b escaped"`,
			expectedResult: "test string \\b escaped",
			err:            false,
		},
		{
			name:           "escape quote err",
			json:           `"test string \n escaped"`,
			expectedResult: "test string \\n escaped",
			err:            false,
		},
		{
			name:           "escape quote err",
			json:           `"test string \\\" escaped`,
			expectedResult: `test string \\\`,
			err:            false,
		},
		{
			name:           "escape quote err",
			json:           `"test string \\\l escaped"`,
			expectedResult: `test string \\\l escaped`,
			err:            false,
		},
		{
			name:           "invalid-json",
			json:           `invalid`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "string-complex",
			json:           `  "string with spaces and \"escape\"d \"quotes\" and escaped line returns \n and escaped \\\\ escaped char"`,
			expectedResult: "string with spaces and \\",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			str := ""
			dec := NewDecoder(strings.NewReader(testCase.json))
			err := dec.StringNoEscape(&str)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should of the given type")
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(
				t,
				testCase.expectedResult,
				str,
				"'%s' should be equal to expectedResult",
				str,
			)
		})
	}
}

func TestDecoderStringInvalidType(t *testing.T) {
	t.Parallel()

	json := []byte(`1`)
	var v string
	err := Unmarshal(json, &v)
	require.Error(t, err, "Err must not be nil as JSON is invalid")
	assert.IsType(t, InvalidUnmarshalError(""), err, "err message must be 'Invalid JSON'")
}

func TestDecoderStringDecoderAPI(t *testing.T) {
	t.Parallel()

	var v string
	dec := NewDecoder(strings.NewReader(`"hello world!"`))
	defer dec.Release()
	err := dec.DecodeString(&v)
	require.NoError(t, err)
	assert.Equal(t, "hello world!", v, "v must be equal to 'hello world!'")
}

func TestDecoderStringPoolError(t *testing.T) {
	t.Parallel()

	// reset the pool to make sure it's not full
	decPool = sync.Pool{
		New: func() any {
			return NewDecoder(nil)
		},
	}
	result := ""
	dec := NewDecoder(nil)
	dec.Release()
	defer func() {
		err := recover()
		require.Error(t, err.(error), "err shouldn't be nil")
		assert.IsType(t, InvalidUsagePooledDecoderError(""), err, "err should be of type InvalidUsagePooledDecoderError")
	}()
	_ = dec.DecodeString(&result)
	assert.True(t, false, "should not be called as decoder should have panicked")
}

func TestDecoderSkipEscapedStringError(t *testing.T) {
	t.Parallel()

	dec := NewDecoder(strings.NewReader(``))
	defer dec.Release()
	err := dec.skipEscapedString()
	require.Error(t, err, "Err must be nil")
	assert.IsType(t, InvalidJSONError(""), err, "err must be of type InvalidJSONError")
}

func TestDecoderSkipEscapedStringError2(t *testing.T) {
	t.Parallel()

	dec := NewDecoder(strings.NewReader(`\"`))
	defer dec.Release()
	err := dec.skipEscapedString()
	require.Error(t, err, "Err must be nil")
	assert.IsType(t, InvalidJSONError(""), err, "err must be of type InvalidJSONError")
}

func TestDecoderSkipEscapedStringError3(t *testing.T) {
	t.Parallel()

	dec := NewDecoder(strings.NewReader(`invalid`))
	defer dec.Release()
	err := dec.skipEscapedString()
	require.Error(t, err, "Err must be nil")
	assert.IsType(t, InvalidJSONError(""), err, "err must be of type InvalidJSONError")
}

func TestDecoderSkipEscapedStringError4(t *testing.T) {
	t.Parallel()

	dec := NewDecoder(strings.NewReader(`\u12`))
	defer dec.Release()
	err := dec.skipEscapedString()
	require.Error(t, err, "Err must be nil")
	assert.IsType(t, InvalidJSONError(""), err, "err must be of type InvalidJSONError")
}

func TestDecoderSkipStringError(t *testing.T) {
	t.Parallel()

	dec := NewDecoder(strings.NewReader(`invalid`))
	defer dec.Release()
	err := dec.skipString()
	require.Error(t, err, "Err must be nil")
	assert.IsType(t, InvalidJSONError(""), err, "err must be of type InvalidJSONError")
}

func TestSkipString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult string
		err            bool
		errType        any
	}{
		{
			name:           "escape quote err",
			json:           `test string \\" escaped"`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "escape quote err",
			json:           `test string \\\l escaped"`,
			expectedResult: ``,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "string-solidus",
			json:           `Asia\/Bangkok","enable":true}"`,
			expectedResult: "",
			err:            false,
		},
		{
			name:           "string-unicode",
			json:           `[2]\u66fe\u5b97\u5357"`,
			expectedResult: "",
			err:            false,
		},
	}

	for _, testCase := range testCases {
		dec := NewDecoder(strings.NewReader(testCase.json))
		err := dec.skipString()
		if testCase.err {
			require.Error(t, err)
			if testCase.errType != nil {
				assert.IsType(
					t,
					testCase.errType,
					err,
					"err should be of expected type",
				)
			}
			return
		}
		require.NoError(t, err)
	}
}
