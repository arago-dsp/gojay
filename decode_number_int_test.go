package gojay

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoderInt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int
		err            bool
		errType        any
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           "1039405",
			expectedResult: 1039405,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-big",
			json:           "9223372036854775807",
			expectedResult: 9223372036854775807,
		},
		{
			name:           "basic-big-overflow",
			json:           "9223372036854775808",
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-big-overflow2",
			json:           "92233720368547758089",
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-big-overflow3",
			json:           "92233720368547758089 ",
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-2349557",
			expectedResult: -2349557,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876 ",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1e2",
			expectedResult: 100,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+06",
			expectedResult: 5000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5.01e+10",
			expectedResult: 50100000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+005",
			expectedResult: 800000,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5e-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+06",
			expectedResult: -5000000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+005",
			expectedResult: -800000,
		},
		{
			name:           "error1",
			json:           "132zz4",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "negative-error2",
			json:           " -1213xdde2323 ",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error3",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},

		{
			name:           "error7",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "error_divide_by_zero",
			json:           "0e-8",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "100000e-4",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			var v int
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil && err != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedResult, v, "v must be equal to %d", testCase.expectedResult)
		})
	}
	t.Run("pool-error", func(t *testing.T) {
		t.Parallel()

		result := 1
		dec := NewDecoder(nil)
		dec.Release()
		defer func() {
			err := recover()
			require.Error(t, err.(error), "err shouldn't be nil")
			assert.IsType(
				t,
				InvalidUsagePooledDecoderError(""),
				err,
				"err should be of type InvalidUsagePooledDecoderError",
			)
		}()
		_ = dec.DecodeInt(&result)
		assert.True(t, false, "should not be called as decoder should have panicked")
	})
	t.Run("decoder-api", func(t *testing.T) {
		t.Parallel()

		var v int
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.DecodeInt(&v)
		require.NoError(t, err)
		assert.Equal(t, int(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api2", func(t *testing.T) {
		t.Parallel()

		var v int
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.Decode(&v)
		require.NoError(t, err)
		assert.Equal(t, 33, v, "v must be equal to 33")
	})
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		var v int
		dec := NewDecoder(strings.NewReader(``))
		defer dec.Release()
		err := dec.DecodeInt(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderIntNull(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int
		err            bool
		errType        any
		resultIsNil    bool
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           "1039405",
			expectedResult: 1039405,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
			resultIsNil:    true,
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-big",
			json:           "9223372036854775807",
			expectedResult: 9223372036854775807,
		},
		{
			name:           "basic-big-overflow",
			json:           "9223372036854775808",
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-big-overflow2",
			json:           "92233720368547758089",
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-big-overflow3",
			json:           "92233720368547758089 ",
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-2349557",
			expectedResult: -2349557,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876 ",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1e2",
			expectedResult: 100,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+06",
			expectedResult: 5000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5.01e+10",
			expectedResult: 50100000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+005",
			expectedResult: 800000,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5e-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+06",
			expectedResult: -5000000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+005",
			expectedResult: -800000,
		},
		{
			name:           "error1",
			json:           "132zz4",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "negative-error2",
			json:           " -1213xdde2323 ",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error3",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},

		{
			name:           "error7",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "error_divide_by_zero",
			json:           "0e-8",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "100000e-4",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			v := (*int)(nil)
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil && err != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s",
						reflect.TypeOf(err).String(),
					)
				}
				return
			}
			require.NoError(t, err)
			if testCase.resultIsNil {
				assert.Nil(t, v)
			} else {
				assert.Equal(
					t,
					testCase.expectedResult,
					*v,
					"v must be equal to %d",
					testCase.expectedResult,
				)
			}
		})
	}
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		v := new(int)
		err := Unmarshal([]byte(``), &v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
	t.Run("decoder-api-invalid-json2", func(t *testing.T) {
		t.Parallel()

		v := new(int)
		dec := NewDecoder(strings.NewReader(``))
		err := dec.IntNull(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int64
		err            bool
		errType        any
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 1039405",
			expectedResult: 1039405,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-big",
			json:           "9223372036854775807",
			expectedResult: 9223372036854775807,
		},
		{
			name:           "basic-big-overflow",
			json:           " 9223372036854775808",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           " 9223372036854775827",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "92233720368547758089",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow3",
			json:           "92233720368547758089 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-negative2",
			json:           "-2349557",
			expectedResult: -2349557,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1e2",
			expectedResult: 100,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+06 ",
			expectedResult: 5000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+005",
			expectedResult: 800000,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5e-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "before-exp-err-too-big",
			json:           "10.11231242345325435464364643e1",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+06",
			expectedResult: -5000000,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5.4e+06",
			expectedResult: -5400000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+005",
			expectedResult: -800000,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e1000000000000000000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e1000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "8ea+00a5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error1",
			json:           "132zz4",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "error_divide_by_zero",
			json:           "0e-8",
			expectedResult: 0,
		},
		{
			name:           "error_divide_by_zero",
			json:           "0e-8",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "100000e-4",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int64 {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int64(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int64 {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int64(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			var v int64
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.expectedResult, v, "v must be equal to %d", testCase.expectedResult)
		})
	}
	t.Run("pool-error", func(t *testing.T) {
		t.Parallel()

		result := int64(1)
		dec := NewDecoder(nil)
		dec.Release()
		defer func() {
			err := recover()
			require.Error(t, err.(error), "err shouldn't be nil")
			assert.IsType(
				t,
				InvalidUsagePooledDecoderError(""),
				err,
				"err should be of type InvalidUsagePooledDecoderError",
			)
		}()
		_ = dec.DecodeInt64(&result)
		assert.True(t, false, "should not be called as decoder should have panicked")
	})
	t.Run("decoder-api", func(t *testing.T) {
		t.Parallel()

		var v int64
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.DecodeInt64(&v)
		require.NoError(t, err)
		assert.Equal(t, int64(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api2", func(t *testing.T) {
		t.Parallel()

		var v int64
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.Decode(&v)
		require.NoError(t, err)
		assert.Equal(t, int64(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		var v int64
		dec := NewDecoder(strings.NewReader(``))
		defer dec.Release()
		err := dec.DecodeInt64(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt64Null(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int64
		err            bool
		errType        any
		resultIsNil    bool
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 1039405",
			expectedResult: 1039405,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
			resultIsNil:    true,
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-big",
			json:           "9223372036854775807",
			expectedResult: 9223372036854775807,
		},
		{
			name:           "basic-big-overflow",
			json:           " 9223372036854775808",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           " 9223372036854775827",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "92233720368547758089",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow3",
			json:           "92233720368547758089 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-negative2",
			json:           "-2349557",
			expectedResult: -2349557,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1e2",
			expectedResult: 100,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+06 ",
			expectedResult: 5000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+005",
			expectedResult: 800000,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5e-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+06",
			expectedResult: -5000000,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5.4e+06",
			expectedResult: -5400000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+005",
			expectedResult: -800000,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e1000000000000000000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e1000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "8ea+00a5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error1",
			json:           "132zz4",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "error_divide_by_zero",
			json:           "0e-8",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "100000e-4",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int64 {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int64(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int64 {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int64(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			v := (*int64)(nil)
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
				return
			}
			require.NoError(t, err)
			if testCase.resultIsNil {
				assert.Nil(t, v)
			} else {
				assert.Equal(
					t,
					testCase.expectedResult,
					*v,
					"v must be equal to %d",
					testCase.expectedResult,
				)
			}
		})
	}
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		v := new(int64)
		err := Unmarshal([]byte(``), &v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
	t.Run("decoder-api-invalid-json2", func(t *testing.T) {
		t.Parallel()

		v := new(int64)
		dec := NewDecoder(strings.NewReader(``))
		err := dec.Int64Null(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt32(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int32
		err            bool
		errType        any
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 1039405",
			expectedResult: 1039405,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-2349557",
			expectedResult: -2349557,
		},
		{
			name:           "basic-big",
			json:           " 2147483647",
			expectedResult: 2147483647,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483648",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483657",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "21474836483",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1.2E2",
			expectedResult: 120,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+005 ",
			expectedResult: 350000,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+005",
			expectedResult: 350000,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+06",
			expectedResult: 5000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+005 ",
			expectedResult: 800000,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5E-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+06",
			expectedResult: -5000000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+005",
			expectedResult: -800000,
		},
		{
			name:           "before-exp-err-too-big",
			json:           "10.11231242345325435464364643e1",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e1000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e1000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e100000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e100000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-float",
			json:           "8.32 ",
			expectedResult: 8,
		},
		{
			name:           "error",
			json:           "83zez4",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error",
			json:           "8ea00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error2",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "error_divide_by_zero",
			json:           "0e-8",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "100000e-4",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int32 {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int32(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int32 {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int32(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			var v int32
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(
				t,
				testCase.expectedResult,
				v,
				"v must be equal to %d", testCase.expectedResult,
			)
		})
	}
	t.Run("pool-error", func(t *testing.T) {
		t.Parallel()

		result := int32(1)
		dec := NewDecoder(nil)
		dec.Release()
		defer func() {
			err := recover()
			require.Error(t, err.(error), "err shouldn't be nil")
			assert.IsType(t, InvalidUsagePooledDecoderError(""), err, "err should be of type InvalidUsagePooledDecoderError")
		}()
		_ = dec.DecodeInt32(&result)
		assert.True(t, false, "should not be called as decoder should have panicked")
	})
	t.Run("decoder-api", func(t *testing.T) {
		t.Parallel()

		var v int32
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.DecodeInt32(&v)
		require.NoError(t, err)
		assert.Equal(t, int32(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api2", func(t *testing.T) {
		t.Parallel()

		var v int32
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.Decode(&v)
		require.NoError(t, err)
		assert.Equal(t, int32(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		var v int32
		dec := NewDecoder(strings.NewReader(``))
		defer dec.Release()
		err := dec.DecodeInt32(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt32Null(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int32
		err            bool
		errType        any
		resultIsNil    bool
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 1039405",
			expectedResult: 1039405,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
			resultIsNil:    true,
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-2349557",
			expectedResult: -2349557,
		},
		{
			name:           "basic-big",
			json:           " 2147483647",
			expectedResult: 2147483647,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483648",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483657",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "21474836483",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1.2E2",
			expectedResult: 120,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+005 ",
			expectedResult: 350000,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+005",
			expectedResult: 350000,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+06",
			expectedResult: 5000000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+005 ",
			expectedResult: 800000,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5E-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+06",
			expectedResult: -5000000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+005",
			expectedResult: -800000,
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e1000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e1000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e100000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e100000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-float",
			json:           "8.32 ",
			expectedResult: 8,
		},
		{
			name:           "error",
			json:           "83zez4",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error",
			json:           "8ea00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error2",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "100000e-4",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int32 {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int32(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int32 {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int32(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			v := (*int32)(nil)
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
				return
			}

			require.NoError(t, err)
			if testCase.resultIsNil {
				assert.Nil(t, v)
			} else {
				assert.Equal(
					t,
					testCase.expectedResult,
					*v,
					"v must be equal to %d",
					testCase.expectedResult,
				)
			}
		})
	}
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		v := new(int32)
		err := Unmarshal([]byte(``), &v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
	t.Run("decoder-api-invalid-json2", func(t *testing.T) {
		t.Parallel()

		v := new(int32)
		dec := NewDecoder(strings.NewReader(``))
		err := dec.Int32Null(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt16(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int16
		err            bool
		errType        any
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 5321",
			expectedResult: 5321,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-2456",
			expectedResult: -2456,
		},
		{
			name:           "basic-big",
			json:           " 24566",
			expectedResult: 24566,
		},
		{
			name:           "basic-big-overflow",
			json:           "66535",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           "32768",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483648",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "21474836483",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1.2E2",
			expectedResult: 120,
		},
		{
			name: "exponent too big",
			json: "1000.202302302422324435342E2",
			err:  true,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+001 ",
			expectedResult: 35,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+002",
			expectedResult: 350,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+03",
			expectedResult: 5000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+02 ",
			expectedResult: 800,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5E-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+03",
			expectedResult: -5000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+003",
			expectedResult: -8000,
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "8.32 ",
			expectedResult: 8,
		},
		{
			name:           "error",
			json:           "83zez4",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error",
			json:           "8ea00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error2",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "0.e",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error8",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "error_divide_by_zero",
			json:           "0e-8",
			expectedResult: 0,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			var v int16
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(
				t,
				testCase.expectedResult,
				v,
				"v must be equal to %d", testCase.expectedResult,
			)
		})
	}
	t.Run("pool-error", func(t *testing.T) {
		t.Parallel()

		result := int16(1)
		dec := NewDecoder(nil)
		dec.Release()
		defer func() {
			err := recover()
			require.Error(t, err.(error), "err shouldn't be nil")
			assert.IsType(t, InvalidUsagePooledDecoderError(""), err, "err should be of type InvalidUsagePooledDecoderError")
		}()
		_ = dec.DecodeInt16(&result)
		assert.True(t, false, "should not be called as decoder should have panicked")
	})
	t.Run("decoder-api", func(t *testing.T) {
		t.Parallel()

		var v int16
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.DecodeInt16(&v)
		require.NoError(t, err)
		assert.Equal(t, int16(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api2", func(t *testing.T) {
		t.Parallel()

		var v int16
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.Decode(&v)
		require.NoError(t, err)
		assert.Equal(t, int16(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		var v int16
		dec := NewDecoder(strings.NewReader(``))
		defer dec.Release()
		err := dec.DecodeInt16(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt16Null(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int16
		err            bool
		errType        any
		resultIsNil    bool
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 5321",
			expectedResult: 5321,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
			resultIsNil:    true,
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-2456",
			expectedResult: -2456,
		},
		{
			name:           "basic-big",
			json:           " 24566",
			expectedResult: 24566,
		},
		{
			name:           "basic-big-overflow",
			json:           "66535",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           "32768",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483648",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "21474836483",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1.2E2",
			expectedResult: 120,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+001 ",
			expectedResult: 35,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+002",
			expectedResult: 350,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+03",
			expectedResult: 5000,
		},
		{
			name:           "basic-exponent-positive-positive-exp3",
			json:           "3e+3",
			expectedResult: 3000,
		},
		{
			name:           "basic-exponent-positive-positive-exp4",
			json:           "8e+02 ",
			expectedResult: 800,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5E-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-005",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0e10000000000 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+03",
			expectedResult: -5000,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e03",
			expectedResult: -3000,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+003",
			expectedResult: -8000,
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "8.32 ",
			expectedResult: 8,
		},
		{
			name:           "error",
			json:           "83zez4",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error",
			json:           "8ea00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error2",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "0.e",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error8",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "10000e-3",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int16 {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int16(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int16 {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int16(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			v := (*int16)(nil)
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
				return
			}
			require.NoError(t, err)
			if testCase.resultIsNil {
				assert.Nil(t, v)
			} else {
				assert.Equal(
					t,
					testCase.expectedResult,
					*v,
					"v must be equal to %d",
					testCase.expectedResult,
				)
			}
		})
	}
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		v := new(int16)
		err := Unmarshal([]byte(``), &v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
	t.Run("decoder-api-invalid-json2", func(t *testing.T) {
		t.Parallel()

		v := new(int16)
		dec := NewDecoder(strings.NewReader(``))
		err := dec.Int16Null(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt8(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int8
		err            bool
		errType        any
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 127",
			expectedResult: 127,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-123",
			expectedResult: -123,
		},
		{
			name:           "basic-big",
			json:           " 43",
			expectedResult: 43,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483648",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           "137",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           "128",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "21474836483",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1.2E2",
			expectedResult: 120,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+001 ",
			expectedResult: 35,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+001",
			expectedResult: 35,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+01",
			expectedResult: 50,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5E-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-1 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e1 ",
			expectedResult: 80,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-1",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+01",
			expectedResult: -50,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e01",
			expectedResult: -30,
		},

		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "before-exp-err-too-big",
			json:           "10.11231242345325435464364643e1",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+001",
			expectedResult: -80,
		},
		{
			name:           "exponent-err-too-big2",
			json:           "0e100 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big2",
			json:           "0.1e100 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "8.32 ",
			expectedResult: 8,
		},
		{
			name:           "error",
			json:           "83zez4",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error",
			json:           "8ea00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error2",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "0.e",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error8",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error8",
			json:           "-5.01e",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-exponent-positive-negative-exp1",
			json:           "100e-1",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int8 {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int8(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int8 {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int8(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			var v int8
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s", reflect.TypeOf(err).String(),
					)
				}
			} else {
				require.NoError(t, err)
			}
			assert.Equal(
				t,
				testCase.expectedResult,
				v,
				"v must be equal to %d",
				testCase.expectedResult,
			)
		})
	}
	t.Run("pool-error", func(t *testing.T) {
		t.Parallel()

		result := int8(1)
		dec := NewDecoder(nil)
		dec.Release()
		defer func() {
			err := recover()
			require.Error(t, err.(error), "err shouldn't be nil")
			assert.IsType(t, InvalidUsagePooledDecoderError(""), err, "err should be of type InvalidUsagePooledDecoderError")
		}()
		_ = dec.DecodeInt8(&result)
		assert.True(t, false, "should not be called as decoder should have panicked")
	})
	t.Run("decoder-api", func(t *testing.T) {
		t.Parallel()
		var v int8
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.DecodeInt8(&v)
		require.NoError(t, err)
		assert.Equal(t, int8(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api2", func(t *testing.T) {
		t.Parallel()

		var v int8
		dec := NewDecoder(strings.NewReader(`33`))
		defer dec.Release()
		err := dec.Decode(&v)
		require.NoError(t, err)
		assert.Equal(t, int8(33), v, "v must be equal to 33")
	})
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		var v int8
		dec := NewDecoder(strings.NewReader(``))
		defer dec.Release()
		err := dec.DecodeInt8(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}

func TestDecoderInt8Null(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult int8
		err            bool
		errType        any
		resultIsNil    bool
	}{
		{
			name:           "basic-positive",
			json:           "100",
			expectedResult: 100,
		},
		{
			name:           "basic-positive2",
			json:           " 127",
			expectedResult: 127,
		},
		{
			name:           "basic-negative",
			json:           "-2",
			expectedResult: -2,
		},
		{
			name:           "basic-null",
			json:           "null",
			expectedResult: 0,
			resultIsNil:    true,
		},
		{
			name:           "basic-negative-err",
			json:           "-",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative-err",
			json:           "-q",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-null-err",
			json:           "nxll",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-skip-data-err",
			json:           "trua",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "basic-negative2",
			json:           "-123",
			expectedResult: -123,
		},
		{
			name:           "basic-big",
			json:           " 43",
			expectedResult: 43,
		},
		{
			name:           "basic-big-overflow",
			json:           " 2147483648",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           "137",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow",
			json:           "128",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-big-overflow2",
			json:           "21474836483",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "2.4595",
			expectedResult: 2,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876",
			expectedResult: -7,
		},
		{
			name:           "basic-float2",
			json:           "-7.8876a",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-positive-positive-exp",
			json:           "1.2E2",
			expectedResult: 120,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+001 ",
			expectedResult: 35,
		},
		{
			name:           "basic-exponent-positive-positive-exp1",
			json:           "3.5e+001",
			expectedResult: 35,
		},
		{
			name:           "basic-exponent-positive-positive-exp2",
			json:           "5e+01",
			expectedResult: 50,
		},
		{
			name:           "basic-exponent-positive-negative-exp",
			json:           "1e-2 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp2",
			json:           "5E-6",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp3",
			json:           "3e-3",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-1 ",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e1 ",
			expectedResult: 80,
		},
		{
			name:           "basic-exponent-positive-negative-exp4",
			json:           "8e-1",
			expectedResult: 0,
		},
		{
			name:           "basic-exponent-negative-positive-exp",
			json:           "-1e2",
			expectedResult: -100,
		},
		{
			name:           "basic-exponent-negative-positive-exp2",
			json:           "-5e+01",
			expectedResult: -50,
		},
		{
			name:           "basic-exponent-negative-positive-exp3",
			json:           "-3e01",
			expectedResult: -30,
		},
		{
			name:           "error3",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "exponent-err-",
			json:           "0.1e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1e10000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big",
			json:           "0.1932242242424244244e1000000000000000000000000",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-negative-positive-exp4",
			json:           "-8e+001",
			expectedResult: -80,
		},
		{
			name:           "exponent-err-too-big2",
			json:           "0e100 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "exponent-err-too-big2",
			json:           "0.1e100 ",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-exponent-err",
			json:           "3e",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "basic-float",
			json:           "8.32 ",
			expectedResult: 8,
		},
		{
			name:           "error",
			json:           "83zez4",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error",
			json:           "8ea00$aa5",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error2",
			json:           "-8e+00$aa5",
			expectedResult: 0,
			err:            true,
		},
		{
			name:           "error4",
			json:           "0.E----",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error5",
			json:           "0E40",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error6",
			json:           "0.e-9",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error7",
			json:           "0.e",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error8",
			json:           "-5.e-2",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "error8",
			json:           "-5.01e",
			expectedResult: 0,
			err:            true,
			errType:        InvalidJSONError(""),
		},
		{
			name:           "invalid-type",
			json:           `"string"`,
			expectedResult: 0,
			err:            true,
			errType:        InvalidUnmarshalError(""),
		},
		{
			name:           "basic-exponent-positive-negative-exp1",
			json:           "100e-1",
			expectedResult: 10,
		},
		{
			name: "fuzz-crasher",
			json: "0E-0",
			expectedResult: func() int8 {
				r, err := strconv.ParseFloat("0E-0", 64)
				require.NoError(t, err)

				return int8(r)
			}(),
			err: false,
		},
		{
			name: "exponent-negative-zero",
			json: "10E-0",
			expectedResult: func() int8 {
				r, err := strconv.ParseFloat("10E-0", 64)
				require.NoError(t, err)

				return int8(r)
			}(),
			err: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			json := []byte(testCase.json)
			v := (*int8)(nil)
			err := Unmarshal(json, &v)
			if testCase.err {
				require.Error(t, err)
				if testCase.errType != nil {
					assert.IsType(
						t,
						testCase.errType,
						err,
						"err should be of type %s",
						reflect.TypeOf(err).String(),
					)
				}
				return
			}
			require.NoError(t, err)
			if testCase.resultIsNil {
				assert.Nil(t, v)
			} else {
				assert.Equal(
					t,
					testCase.expectedResult,
					*v,
					"v must be equal to %d",
					testCase.expectedResult,
				)
			}
		})
	}
	t.Run("decoder-api-invalid-json", func(t *testing.T) {
		t.Parallel()

		v := new(int8)
		err := Unmarshal([]byte(``), &v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
	t.Run("decoder-api-invalid-json2", func(t *testing.T) {
		t.Parallel()

		v := new(int8)
		dec := NewDecoder(strings.NewReader(``))
		err := dec.Int8Null(&v)
		require.Error(t, err)
		assert.IsType(t, InvalidJSONError(""), err, "err should be of type InvalidJSONError")
	})
}
