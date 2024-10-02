package gojay

// AddSliceString unmarshal the next JSON array of strings to the given *[]string s.
func (dec *Decoder) AddSliceString(s *[]string) error {
	return dec.SliceString(s)
}

// SliceString unmarshal the next JSON array of strings to the given *[]string s.
func (dec *Decoder) SliceString(s *[]string) error {
	str := AcquireString()
	defer ReleaseString(str)

	return dec.Array(DecodeArrayFunc(func(dec *Decoder) error {
		if err := dec.String(str); err != nil {
			return err
		}
		*s = append(*s, *str)
		return nil
	}))
}

// AddSliceInt unmarshal the next JSON array of integers to the given *[]int s.
func (dec *Decoder) AddSliceInt(s *[]int) error {
	return dec.SliceInt(s)
}

// SliceInt unmarshal the next JSON array of integers to the given *[]int s.
func (dec *Decoder) SliceInt(s *[]int) error {
	var i int
	return dec.Array(DecodeArrayFunc(func(dec *Decoder) error {
		if err := dec.Int(&i); err != nil {
			return err
		}
		*s = append(*s, i)
		return nil
	}))
}

// AddSliceInt8 unmarshal the next JSON array of integers to the given *[]int s.
func (dec *Decoder) AddSliceInt8(s *[]int8) error {
	return dec.SliceInt8(s)
}

// SliceInt8 unmarshal the next JSON array of integers to the given *[]int8 s.
func (dec *Decoder) SliceInt8(s *[]int8) error {
	var i int8
	return dec.Array(DecodeArrayFunc(func(dec *Decoder) error {
		if err := dec.Int8(&i); err != nil {
			return err
		}
		*s = append(*s, i)
		return nil
	}))
}

// AddSliceUint8 unmarshal the next JSON array of integers to the given *[]uint8 s.
func (dec *Decoder) AddSliceUint8(s *[]uint8) error {
	return dec.SliceUint8(s)
}

// SliceUint8 unmarshal the next JSON array of integers to the given *[]uint8 s.
func (dec *Decoder) SliceUint8(s *[]uint8) error {
	var i uint8
	return dec.Array(DecodeArrayFunc(func(dec *Decoder) error {
		if err := dec.Uint8(&i); err != nil {
			return err
		}
		*s = append(*s, i)
		return nil
	}))
}

// AddSliceFloat64 unmarshal the next JSON array of floats to the given *[]float64 s.
func (dec *Decoder) AddSliceFloat64(s *[]float64) error {
	return dec.SliceFloat64(s)
}

// SliceFloat64 unmarshal the next JSON array of floats to the given *[]float64 s.
func (dec *Decoder) SliceFloat64(s *[]float64) error {
	var i float64
	return dec.Array(DecodeArrayFunc(func(dec *Decoder) error {
		if err := dec.Float64(&i); err != nil {
			return err
		}
		*s = append(*s, i)
		return nil
	}))
}

// AddSliceBool unmarshal the next JSON array of bool to the given *[]bool s.
func (dec *Decoder) AddSliceBool(s *[]bool) error {
	return dec.SliceBool(s)
}

// SliceBool unmarshal the next JSON array of bool to the given *[]bool s.
func (dec *Decoder) SliceBool(s *[]bool) error {
	var b bool
	return dec.Array(DecodeArrayFunc(func(dec *Decoder) error {
		if err := dec.Bool(&b); err != nil {
			return err
		}
		*s = append(*s, b)
		return nil
	}))
}

// AddSliceStringNoEscape unmarshal the next JSON array of strings to the given *[]string s.
func (dec *Decoder) AddSliceStringNoEscape(s *[]string) error {
	return dec.SliceStringNoEscape(s)
}

// SliceStringNoEscape unmarshal the next JSON array of strings to the given *[]string s.
func (dec *Decoder) SliceStringNoEscape(s *[]string) error {
	str := AcquireString()
	defer ReleaseString(str)

	return dec.Array(DecodeArrayFunc(func(dec *Decoder) error {
		if err := dec.StringNoEscape(str); err != nil {
			return err
		}
		*s = append(*s, *str)
		return nil
	}))
}
