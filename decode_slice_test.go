package gojay

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type slicesTestObject struct {
	sliceString         []string
	sliceStringNoEscape []string
	sliceInt            []int
	sliceInt8           []int8
	sliceFloat64        []float64
	sliceBool           []bool
}

func (s *slicesTestObject) UnmarshalJSONObject(dec *Decoder, k string) error {
	switch k {
	case "sliceString":
		return dec.AddSliceString(&s.sliceString)
	case "sliceStringNoEscape":
		return dec.AddSliceStringNoEscape(&s.sliceStringNoEscape)
	case "sliceInt":
		return dec.AddSliceInt(&s.sliceInt)
	case "sliceInt8":
		return dec.AddSliceInt8(&s.sliceInt8)
	case "sliceFloat64":
		return dec.AddSliceFloat64(&s.sliceFloat64)
	case "sliceBool":
		return dec.AddSliceBool(&s.sliceBool)
	}
	return nil
}

func (s *slicesTestObject) NKeys() int {
	return 4
}

func TestDecodeSlices(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		json           string
		expectedResult slicesTestObject
		err            bool
	}{
		{
			name: "basic slice string",
			json: `{
				"sliceString": ["foo","bar"]
			}`,
			expectedResult: slicesTestObject{
				sliceString: []string{"foo", "bar"},
			},
		},
		{
			name: "basic slice string no escape",
			json: `{
				"sliceStringNoEscape": ["foo","bar"]
			}`,
			expectedResult: slicesTestObject{
				sliceStringNoEscape: []string{"foo", "bar"},
			},
		},
		{
			name: "basic slice bool",
			json: `{
				"sliceBool": [true,false]
			}`,
			expectedResult: slicesTestObject{
				sliceBool: []bool{true, false},
			},
		},
		{
			name: "basic slice int",
			json: `{
				"sliceInt": [1,2,3]
			}`,
			expectedResult: slicesTestObject{
				sliceInt: []int{1, 2, 3},
			},
		},
		{
			name: "basic slice int8",
			json: `{
				"sliceInt8": [1,2,3]
			}`,
			expectedResult: slicesTestObject{
				sliceInt8: []int8{1, 2, 3},
			},
		},
		{
			name: "basic slice float64",
			json: `{
				"sliceFloat64": [1.3,2.4,3.1]
			}`,
			expectedResult: slicesTestObject{
				sliceFloat64: []float64{1.3, 2.4, 3.1},
			},
		},
		{
			name: "err slice float64",
			json: `{
				"sliceFloat64": [1.3",2.4,3.1]
			}`,
			err: true,
		},
		{
			name: "err slice str",
			json: `{
				"sliceString": [",""]
			}`,
			err: true,
		},
		{
			name: "err slice str no escape",
			json: `{
				"sliceStringNoEscape": [",""]
			}`,
			err: true,
		},
		{
			name: "err slice int",
			json: `{
				"sliceInt": [1t,2,3]
			}`,
			err: true,
		},
		{
			name: "err slice int8",
			json: `{
				"sliceInt8": [1t,2,3]
			}`,
			err: true,
		},
		{
			name: "err slice bool",
			json: `{
				"sliceBool": [truo,false]
			}`,
			err: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				t.Parallel()

				dec := BorrowDecoder(strings.NewReader(testCase.json))
				var o slicesTestObject
				err := dec.Decode(&o)

				if testCase.err {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, o)
			},
		)
	}
}
