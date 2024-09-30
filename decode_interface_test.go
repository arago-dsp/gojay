package gojay

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeInterfaceBasic(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		json            string
		expectedResult  any
		err             bool
		errType         any
		skipCheckResult bool
	}{
		{
			name:           "array",
			json:           `[1,2,3]`,
			expectedResult: []any{float64(1), float64(2), float64(3)},
			err:            false,
		},
		{
			name:           "object",
			json:           `{"testStr": "hello world!"}`,
			expectedResult: map[string]any{"testStr": "hello world!"},
			err:            false,
		},
		{
			name:           "string",
			json:           `"hola amigos!"`,
			expectedResult: any("hola amigos!"),
			err:            false,
		},
		{
			name:           "bool-true",
			json:           `true`,
			expectedResult: any(true),
			err:            false,
		},
		{
			name:           "bool-false",
			json:           `false`,
			expectedResult: any(false),
			err:            false,
		},
		{
			name:           "null",
			json:           `null`,
			expectedResult: any(nil),
			err:            false,
		},
		{
			name:           "number",
			json:           `1234`,
			expectedResult: any(float64(1234)),
			err:            false,
		},
		{
			name:            "array-error",
			json:            `["h""o","l","a"]`,
			err:             true,
			errType:         &json.SyntaxError{},
			skipCheckResult: true,
		},
		{
			name:            "object-error",
			json:            `{"testStr" "hello world!"}`,
			err:             true,
			errType:         &json.SyntaxError{},
			skipCheckResult: true,
		},
		{
			name:            "string-error",
			json:            `"hola amigos!`,
			err:             true,
			errType:         InvalidJSONError(""),
			skipCheckResult: true,
		},
		{
			name:            "bool-true-error",
			json:            `truee`,
			err:             true,
			errType:         InvalidJSONError(""),
			skipCheckResult: true,
		},
		{
			name:            "bool-false-error",
			json:            `fase`,
			expectedResult:  any(false),
			err:             true,
			errType:         InvalidJSONError(""),
			skipCheckResult: true,
		},
		{
			name:            "null-error",
			json:            `nulllll`,
			err:             true,
			errType:         InvalidJSONError(""),
			skipCheckResult: true,
		},
		{
			name:            "number-error",
			json:            `1234"`,
			err:             true,
			errType:         InvalidJSONError(""),
			skipCheckResult: true,
		},
		{
			name:            "unknown-error",
			json:            `?`,
			err:             true,
			errType:         InvalidJSONError(""),
			skipCheckResult: true,
		},
		{
			name:            "empty-json-error",
			json:            ``,
			err:             true,
			errType:         InvalidJSONError(""),
			skipCheckResult: true,
		},
	}

	for _, testCase := range testCases {
		t.Run("DecodeInterface()"+testCase.name, func(t *testing.T) {
			var i any
			dec := BorrowDecoder(strings.NewReader(testCase.json))
			defer dec.Release()
			err := dec.DecodeInterface(&i)
			if testCase.err {
				t.Log(err)
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should be of the given type")
				}
				return
			}
			require.NoError(t, err)
			if !testCase.skipCheckResult {
				assert.Equal(t, testCase.expectedResult, i, "value at given index should be the same as expected results")
			}
		})
	}

	for _, testCase := range testCases {
		t.Run("Decode()"+testCase.name, func(t *testing.T) {
			t.Parallel()

			var i any
			dec := BorrowDecoder(strings.NewReader(testCase.json))
			defer dec.Release()
			err := dec.Decode(&i)
			if testCase.err {
				t.Log(err)
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should be of the given type")
				}
				return
			}
			require.NoError(t, err)
			if !testCase.skipCheckResult {
				assert.Equal(t, testCase.expectedResult, i, "value at given index should be the same as expected results")
			}
		})
	}
}

func TestDecodeInterfaceAsInterface(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		json            string
		expectedResult  any
		err             bool
		errType         any
		skipCheckResult bool
	}{
		{
			name: "basic-array",
			json: `{
        "testStr": "hola",
        "testInterface": ["h","o","l","a"]
      }`,
			expectedResult: map[string]any{
				"testStr":       "hola",
				"testInterface": []any{"h", "o", "l", "a"},
			},
			err: false,
		},
		{
			name: "basic-string",
			json: `{
        "testInterface": "漢字"
      }`,
			expectedResult: map[string]any{
				"testInterface": "漢字",
			},
			err: false,
		},
		{
			name: "basic-error",
			json: `{
        "testInterface": ["a""d","i","o","s"]
      }`,
			err:             true,
			errType:         &json.SyntaxError{},
			skipCheckResult: true,
		},
		{
			name: "basic-interface",
			json: `{
        "testInterface": {
          "string": "prost"
        }
      }`,
			expectedResult: map[string]any{
				"testInterface": map[string]any{"string": "prost"},
			},
			err: false,
		},
		{
			name: "complex-interface",
			json: `{
        "testInterface": {
          "number": 1988,
          "string": "prost",
          "array": ["h","o","l","a"],
          "object": {
            "k": "v",
            "a": [1,2,3]
          },
          "array-of-objects": [
            {"k": "v"},
            {"a": "b"}
          ]
        }
      }`,
			expectedResult: map[string]any{
				"testInterface": map[string]any{
					"array-of-objects": []any{
						map[string]any{"k": "v"},
						map[string]any{"a": "b"},
					},
					"number": float64(1988),
					"string": "prost",
					"array":  []any{"h", "o", "l", "a"},
					"object": map[string]any{
						"k": "v",
						"a": []any{float64(1), float64(2), float64(3)},
					},
				},
			},
			err: false,
		},
	}

	for _, testCase := range testCases {
		t.Run("Decode()"+testCase.name, func(t *testing.T) {
			var s any
			dec := BorrowDecoder(strings.NewReader(testCase.json))
			defer dec.Release()
			err := dec.Decode(&s)
			if testCase.err {
				t.Log(err)
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should be of the given type")
				}
				return
			}
			require.NoError(t, err)
			if !testCase.skipCheckResult {
				assert.Equal(t, testCase.expectedResult, s, "value at given index should be the same as expected results")
			}
		})
	}

	for _, testCase := range testCases {
		t.Run("DecodeInterface()"+testCase.name, func(t *testing.T) {
			t.Parallel()

			var s any
			dec := BorrowDecoder(strings.NewReader(testCase.json))
			defer dec.Release()
			err := dec.DecodeInterface(&s)
			if testCase.err {
				t.Log(err)
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should be of the given type")
				}
				return
			}
			require.NoError(t, err)
			if !testCase.skipCheckResult {
				assert.Equal(t, testCase.expectedResult, s, "value at given index should be the same as expected results")
			}
		})
	}
}

func TestDecodeAsTestObject(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		json            string
		expectedResult  testObject
		err             bool
		errType         any
		skipCheckResult bool
	}{
		{
			name: "basic-array",
			json: `{
        "testStr": "hola",
        "testInterface": ["h","o","l","a"]
      }`,
			expectedResult: testObject{
				testStr:       "hola",
				testInterface: []any{"h", "o", "l", "a"},
			},
			err: false,
		},
		{
			name: "basic-string",
			json: `{
        "testInterface": "漢字"
      }`,
			expectedResult: testObject{
				testInterface: any("漢字"),
			},
			err: false,
		},
		{
			name: "basic-error",
			json: `{
        "testInterface": ["a""d","i","o","s"]
      }`,
			err:             true,
			errType:         &json.SyntaxError{},
			skipCheckResult: true,
		},
		{
			name: "mull-interface",
			json: `{
        "testInterface": null,
        "testStr": "adios"
      }`,
			expectedResult: testObject{
				testInterface: any(nil),
				testStr:       "adios",
			},
			err: false,
		},
		{
			name: "basic-interface",
			json: `{
        "testInterface": {
          "string": "prost"
        },
      }`,
			expectedResult: testObject{
				testInterface: map[string]any{"string": "prost"},
			},
			err: false,
		},
		{
			name: "complex-interface",
			json: `{
        "testInterface": {
          "number": 1988,
          "string": "prost",
          "array": ["h","o","l","a"],
          "object": {
            "k": "v",
            "a": [1,2,3]
          },
          "array-of-objects": [
            {"k": "v"},
            {"a": "b"}
          ]
        },
      }`,
			expectedResult: testObject{
				testInterface: map[string]any{
					"array-of-objects": []any{
						map[string]any{"k": "v"},
						map[string]any{"a": "b"},
					},
					"number": float64(1988),
					"string": "prost",
					"array":  []any{"h", "o", "l", "a"},
					"object": map[string]any{
						"k": "v",
						"a": []any{float64(1), float64(2), float64(3)},
					},
				},
			},
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			s := testObject{}
			dec := BorrowDecoder(strings.NewReader(testCase.json))
			defer dec.Release()
			err := dec.Decode(&s)
			if testCase.err {
				t.Log(err)
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(t, testCase.errType, err, "err should be of the given type")
				}
				return
			}
			require.NoError(t, err)
			if !testCase.skipCheckResult {
				assert.Equal(t, testCase.expectedResult, s, "value at given index should be the same as expected results")
			}
		})
	}
}

func TestUnmarshalInterface(t *testing.T) {
	t.Parallel()

	bytes := []byte(`{
    "testInterface": {
      "number": 1988,
      "null": null,
      "string": "prost",
      "array": ["h","o","l","a"],
      "object": {
        "k": "v",
        "a": [1,2,3]
      },
      "array-of-objects": [
        {"k": "v"},
        {"a": "b"}
      ]
    }
	}`)
	v := &testObject{}
	err := Unmarshal(bytes, v)
	require.NoError(t, err)
	expectedInterface := map[string]any{
		"array-of-objects": []any{
			map[string]any{"k": "v"},
			map[string]any{"a": "b"},
		},
		"number": float64(1988),
		"string": "prost",
		"null":   any(nil),
		"array":  []any{"h", "o", "l", "a"},
		"object": map[string]any{
			"k": "v",
			"a": []any{float64(1), float64(2), float64(3)},
		},
	}
	assert.Equal(t, expectedInterface, v.testInterface, "v.testInterface must be equal to the expected one")
}

func TestUnmarshalInterfaceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		json []byte
	}{
		{
			name: "basic",
			json: []byte(`{"testInterface": {"number": 1bc4}}`),
		},
		{
			name: "syntax",
			json: []byte(`{
        "testInterface": {
          "array?": [1,"a", ?]
        }
      }`),
		},
		{
			name: "complex",
			json: []byte(`{
        "testInterface": {
          "number": 1988,
          "string": "prost",
          "array": ["h""o","l","a"],
          "object": {
            "k": "v",
            "a": [1,2,3]
          },
          "array-of-objects": [
            {"k": "v"},
            {"a": "b"}
          ]
        }
      }`),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			v := &testObject{}
			err := Unmarshal(testCase.json, v)
			require.Error(t, err)
			t.Log(err)
			assert.IsType(t, &json.SyntaxError{}, err, "err should be a json.SyntaxError{}")
		})
	}
}

func TestDecodeInterfacePoolError(t *testing.T) {
	t.Parallel()

	result := any(1)
	dec := NewDecoder(nil)
	dec.Release()
	defer func() {
		err := recover()
		require.Error(t, err.(error), "err shouldn't be nil")
		assert.IsType(t, InvalidUsagePooledDecoderError(""), err, "err should be of type InvalidUsagePooledDecoderError")
	}()
	_ = dec.DecodeInterface(&result)
	assert.True(t, false, "should not be called as decoder should have panicked")
}

func TestDecodeNull(t *testing.T) {
	t.Parallel()

	var i any
	dec := BorrowDecoder(strings.NewReader("null"))
	defer dec.Release()
	err := dec.DecodeInterface(&i)
	require.NoError(t, err)
	assert.Equal(t, any(nil), i, "value at given index should be the same as expected results")
}
