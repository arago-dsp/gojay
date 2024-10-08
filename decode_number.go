package gojay

import (
	"math"
)

var digits []int8

const (
	maxInt64toMultiply  = math.MaxInt64 / 10
	maxInt32toMultiply  = math.MaxInt32 / 10
	maxInt16toMultiply  = math.MaxInt16 / 10
	maxInt8toMultiply   = math.MaxInt8 / 10
	maxUint8toMultiply  = math.MaxUint8 / 10
	maxUint16toMultiply = math.MaxUint16 / 10
	maxUint32toMultiply = math.MaxUint32 / 10
	maxUint64toMultiply = math.MaxUint64 / 10
	maxUint32Length     = 10
	maxUint64Length     = 20
	maxUint16Length     = 5
	maxUint8Length      = 3
	maxInt32Length      = 10
	maxInt64Length      = 19
	maxInt16Length      = 5
	maxInt8Length       = 3
	invalidNumber       = int8(-1)
)

var pow10uint64 = [21]uint64{
	0,
	1,
	10,
	100,
	1000,
	10000,
	100000,
	1000000,
	10000000,
	100000000,
	1000000000,
	10000000000,
	100000000000,
	1000000000000,
	10000000000000,
	100000000000000,
	1000000000000000,
	10000000000000000,
	100000000000000000,
	1000000000000000000,
	10000000000000000000,
}

var skipNumberEndCursorIncrement [256]int

func init() {
	digits = make([]int8, 256)
	for i := 0; i < len(digits); i++ {
		digits[i] = invalidNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		digits[i] = i - int8('0')
	}

	for i := range 256 {
		switch i {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E', '+', '-':
			skipNumberEndCursorIncrement[i] = 1
		}
	}
}

func (dec *Decoder) skipNumber() (int, error) {
	end := dec.cursor + 1
	// look for following numbers
	for j := dec.cursor + 1; j < dec.length || dec.read(); j++ {
		end += skipNumberEndCursorIncrement[dec.data[j]]

		switch dec.data[j] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E', '+', '-', ' ', '\n', '\t', '\r':
			continue
		case ',', '}', ']':
			return end, nil
		default:
			// invalid json we expect numbers, dot (single one), comma, or spaces
			return end, dec.raiseInvalidJSONErr(dec.cursor)
		}
	}

	return end, nil
}

func (dec *Decoder) getExponent() (int64, error) {
	start := dec.cursor
	end := dec.cursor
	for ; dec.cursor < dec.length || dec.read(); dec.cursor++ {
		switch dec.data[dec.cursor] { // is positive
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			end = dec.cursor + 1
		case '-':
			dec.cursor++
			exp, err := dec.getExponent()
			return -exp, err
		case '+':
			dec.cursor++
			return dec.getExponent()
		default:
			// if nothing return 0
			// could raise error
			if start == end {
				return 0, dec.raiseInvalidJSONErr(dec.cursor)
			}
			return dec.atoi64(start, end-1), nil
		}
	}
	if start == end {
		return 0, dec.raiseInvalidJSONErr(dec.cursor)
	}
	return dec.atoi64(start, end-1), nil
}
