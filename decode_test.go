package gojay

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDecodeObj struct {
	test string
}

func (t *testDecodeObj) UnmarshalJSONObject(dec *Decoder, key string) error {
	if key == "test" {
		return dec.AddString(&t.test)
	}
	return nil
}

func (t *testDecodeObj) NKeys() int {
	return 1
}

type testDecodeSlice []*testDecodeObj

func (t *testDecodeSlice) UnmarshalJSONArray(dec *Decoder) error {
	obj := &testDecodeObj{}
	if err := dec.AddObject(obj); err != nil {
		return err
	}
	*t = append(*t, obj)
	return nil
}

type allTypeDecodeTestCase struct {
	name         string
	v            any
	d            []byte
	expectations func(err error, v any, t *testing.T)
}

func allTypesTestCases() []allTypeDecodeTestCase {
	return []allTypeDecodeTestCase{
		{
			v:    new(string),
			d:    []byte(`"test string"`),
			name: "test decode string",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*string)
				require.NoError(t, err)
				assert.Equal(t, "test string", *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*string),
			d:    []byte(`"test string"`),
			name: "test decode string",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**string)
				require.NoError(t, err)
				assert.Equal(t, "test string", **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(string),
			d:    []byte(`null`),
			name: "test decode string null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*string)
				require.NoError(t, err)
				assert.Equal(t, "", *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*string),
			d:    []byte(`null`),
			name: "test decode string null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**string)
				require.NoError(t, err)
				assert.Nil(t, *vt, "v must be nil")
			},
		},
		{
			v:    new(*string),
			d:    []byte(`1`),
			name: "test decode string null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**string)
				require.Error(t, err, "err must be nil")
				assert.Nil(t, *vt, "v must be nil")
			},
		},
		{
			v:    new(int),
			d:    []byte(`1`),
			name: "test decode int",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*int)
				require.NoError(t, err)
				assert.Equal(t, 1, *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*int),
			d:    []byte(`1`),
			name: "test decode int",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**int)
				require.NoError(t, err)
				assert.Equal(t, 1, **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*int),
			d:    []byte(`""`),
			name: "test decode int",
			expectations: func(err error, _ any, t *testing.T) {
				require.Error(t, err, "err must be nil")
			},
		},
		{
			v:    new(*int8),
			d:    []byte(`1`),
			name: "test decode int",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**int8)
				require.NoError(t, err)
				assert.Equal(t, int8(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*int8),
			d:    []byte(`""`),
			name: "test decode int",
			expectations: func(err error, _ any, t *testing.T) {
				require.Error(t, err, "err must be nil")
			},
		},
		{
			v:    new(*int16),
			d:    []byte(`1`),
			name: "test decode int",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**int16)
				require.NoError(t, err)
				assert.Equal(t, int16(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*int16),
			d:    []byte(`""`),
			name: "test decode int",
			expectations: func(err error, _ any, t *testing.T) {
				require.Error(t, err, "err must be nil")
			},
		},
		{
			v:    new(int64),
			d:    []byte(`1`),
			name: "test decode int64",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*int64)
				require.NoError(t, err)
				assert.Equal(t, int64(1), *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*int64),
			d:    []byte(`1`),
			name: "test decode int64",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**int64)
				require.NoError(t, err)
				assert.Equal(t, int64(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*int64),
			d:    []byte(`""`),
			name: "test decode int64",
			expectations: func(err error, _ any, t *testing.T) {
				require.Error(t, err, "err must be nil")
			},
		},
		{
			v:    new(uint64),
			d:    []byte(`1`),
			name: "test decode uint64",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*uint64)
				require.NoError(t, err)
				assert.Equal(t, uint64(1), *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*uint64),
			d:    []byte(`1`),
			name: "test decode uint64",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**uint64)
				require.NoError(t, err)
				assert.Equal(t, uint64(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(any),
			d:    []byte(`[{"test":"test"},{"test":"test2"}]`),
			name: "test decode interface",
			expectations: func(err error, v any, t *testing.T) {
				require.NoError(t, err)
				// v is a pointer to an any, we need to extract the content
				vCont := reflect.ValueOf(v).Elem().Interface()
				vt := vCont.([]any)
				assert.Len(t, vt, 2, "len of vt must be 2")
				vt1 := vt[0].(map[string]any)
				assert.Equal(t, "test", vt1["test"], "vt1['test'] must be equal to 'test'")
				vt2 := vt[1].(map[string]any)
				assert.Equal(t, "test2", vt2["test"], "vt2['test'] must be equal to 'test2'")
			},
		},
		{
			v:    new(uint64),
			d:    []byte(`-1`),
			name: "test decode uint64 negative",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*uint64)
				require.Error(t, err, "err must not be nil")
				assert.Equal(t, uint64(0), *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(int32),
			d:    []byte(`1`),
			name: "test decode int32",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*int32)
				require.NoError(t, err)
				assert.Equal(t, int32(1), *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*int32),
			d:    []byte(`1`),
			name: "test decode int32",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**int32)
				require.NoError(t, err)
				assert.Equal(t, int32(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(uint32),
			d:    []byte(`1`),
			name: "test decode uint32",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*uint32)
				require.NoError(t, err)
				assert.Equal(t, uint32(1), *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*uint32),
			d:    []byte(`1`),
			name: "test decode uint32",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**uint32)
				require.NoError(t, err)
				assert.Equal(t, uint32(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(uint32),
			d:    []byte(`-1`),
			name: "test decode uint32 negative",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*uint32)
				require.Error(t, err, "err must not be nil")
				assert.Equal(t, uint32(0), *vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*uint16),
			d:    []byte(`1`),
			name: "test decode uint16",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**uint16)
				require.NoError(t, err)
				assert.Equal(t, uint16(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(*uint8),
			d:    []byte(`1`),
			name: "test decode uint8",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**uint8)
				require.NoError(t, err)
				assert.Equal(t, uint8(1), **vt, "v must be equal to 1")
			},
		},
		{
			v:    new(float64),
			d:    []byte(`1.15`),
			name: "test decode float64",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*float64)
				require.NoError(t, err)
				assert.InDelta(t, float64(1.15), *vt, 0)
			},
		},
		{
			v:    new(*float64),
			d:    []byte(`1.15`),
			name: "test decode float64",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**float64)
				require.NoError(t, err)
				assert.InDelta(t, float64(1.15), **vt, 0)
			},
		},
		{
			v:    new(float64),
			d:    []byte(`null`),
			name: "test decode float64 null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*float64)
				require.NoError(t, err)
				assert.InDelta(t, float64(0), *vt, 0)
			},
		},
		{
			v:    new(*float32),
			d:    []byte(`1.15`),
			name: "test decode float64 null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**float32)
				require.NoError(t, err)
				assert.InDeltaf(t, float32(1.15), **vt, 0, "v must be equal to 1")
			},
		},
		{
			v:    new(bool),
			d:    []byte(`true`),
			name: "test decode bool true",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*bool)
				require.NoError(t, err)
				assert.True(t, *vt)
			},
		},
		{
			v:    new(*bool),
			d:    []byte(`true`),
			name: "test decode bool true",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(**bool)
				require.NoError(t, err)
				assert.True(t, **vt)
			},
		},
		{
			v:    new(bool),
			d:    []byte(`false`),
			name: "test decode bool false",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*bool)
				require.NoError(t, err)
				assert.False(t, *vt)
			},
		},
		{
			v:    new(bool),
			d:    []byte(`null`),
			name: "test decode bool null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*bool)
				require.NoError(t, err)
				assert.False(t, *vt)
			},
		},
		{
			v:    new(testDecodeObj),
			d:    []byte(`{"test":"test"}`),
			name: "test decode object",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*testDecodeObj)
				require.NoError(t, err)
				assert.Equal(t, "test", vt.test, "v.test must be equal to 'test'")
			},
		},
		{
			v:    new(testDecodeObj),
			d:    []byte(`{"test":null}`),
			name: "test decode object null key",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*testDecodeObj)
				require.NoError(t, err)
				assert.Equal(t, "", vt.test, "v.test must be equal to 'test'")
			},
		},
		{
			v:    new(testDecodeObj),
			d:    []byte(`null`),
			name: "test decode object null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*testDecodeObj)
				require.NoError(t, err)
				assert.Equal(t, "", vt.test, "v.test must be equal to 'test'")
			},
		},
		{
			v:    new(testDecodeSlice),
			d:    []byte(`[{"test":"test"}]`),
			name: "test decode slice",
			expectations: func(err error, v any, t *testing.T) {
				vtPtr := v.(*testDecodeSlice)
				vt := *vtPtr
				require.NoError(t, err)
				assert.Len(t, vt, 1, "len of vt must be 1")
				assert.Equal(t, "test", vt[0].test, "vt[0].test must be equal to 'test'")
			},
		},
		{
			v:    new(testDecodeSlice),
			d:    []byte(`[{"test":"test"},{"test":"test2"}]`),
			name: "test decode slice",
			expectations: func(err error, v any, t *testing.T) {
				vtPtr := v.(*testDecodeSlice)
				vt := *vtPtr
				require.NoError(t, err)
				assert.Len(t, vt, 2, "len of vt must be 2")
				assert.Equal(t, "test", vt[0].test, "vt[0].test must be equal to 'test'")
				assert.Equal(t, "test2", vt[1].test, "vt[1].test must be equal to 'test2'")
			},
		},
		{
			v:    new(struct{}),
			d:    []byte(`{"test":"test"}`),
			name: "test decode invalid type",
			expectations: func(err error, v any, t *testing.T) {
				require.Error(t, err, "err must not be nil")
				assert.IsType(t, InvalidUnmarshalError(""), err, "err must be of type InvalidUnmarshalError")
				assert.Equal(t, fmt.Sprintf(invalidUnmarshalErrorMsg, v), err.Error(), "err message should be equal to invalidUnmarshalErrorMsg")
			},
		},
	}
}

// Unmarshal tests.
func TestUnmarshalAllTypes(t *testing.T) {
	t.Parallel()

	for _, testCase := range allTypesTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := Unmarshal(testCase.d, testCase.v)
			testCase.expectations(err, testCase.v, t)
		})
	}
}

// Decode tests.
func TestDecodeAllTypes(t *testing.T) {
	t.Parallel()

	for _, testCase := range allTypesTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			dec := NewDecoder(bytes.NewReader(testCase.d))
			err := dec.Decode(testCase.v)
			testCase.expectations(err, testCase.v, t)
		})
	}
}

func TestUnmarshalJSONObjects(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		v            UnmarshalerJSONObject
		d            []byte
		expectations func(err error, v any, t *testing.T)
	}{
		{
			v:    new(testDecodeObj),
			d:    []byte(`{"test":"test"}`),
			name: "test decode object",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*testDecodeObj)
				require.NoError(t, err)
				assert.Equal(t, "test", vt.test, "v.test must be equal to 'test'")
			},
		},
		{
			v:    new(testDecodeObj),
			d:    []byte(`{"test":null}`),
			name: "test decode object null key",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*testDecodeObj)
				require.NoError(t, err)
				assert.Equal(t, "", vt.test, "v.test must be equal to 'test'")
			},
		},
		{
			v:    new(testDecodeObj),
			d:    []byte(`null`),
			name: "test decode object null",
			expectations: func(err error, v any, t *testing.T) {
				vt := v.(*testDecodeObj)
				require.NoError(t, err)
				assert.Equal(t, "", vt.test, "v.test must be equal to 'test'")
			},
		},
		{
			v:    new(testDecodeObj),
			d:    []byte(`invalid json`),
			name: "test decode object null",
			expectations: func(err error, _ any, t *testing.T) {
				require.Error(t, err, "err must not be nil")
				assert.IsType(t, InvalidJSONError(""), err, "err must be of type InvalidJSONError")
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := UnmarshalJSONObject(testCase.d, testCase.v)
			testCase.expectations(err, testCase.v, t)
		})
	}
}
