package gojay

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type user struct {
	ID      int     `json:"id"`
	Created uint64  `json:"created"`
	Age     float64 `json:"age"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
}

func (u *user) UnmarshalJSONObject(dec *Decoder, key string) error {
	switch key {
	case "id":
		return dec.Int(&u.ID)
	case "created":
		return dec.Uint64(&u.Created)
	case "age":
		return dec.Float64(&u.Age)
	case "name":
		return dec.String(&u.Name)
	case "email":
		return dec.StringNoEscape(&u.Email)
	}
	return nil
}

func (u *user) NKeys() int {
	return 6
}

func TestMarshal(t *testing.T) {
	assert.JSONEq(t,
		`{"age":1, "created":1, "email":"foo@example.com", "id":1, "name":"foo"}`,
		string(marshalHelper(t, 1, 1, 1, "foo", "foo@example.com")))
}

func FuzzUnmarshalRaw(f *testing.F) {
	f.Add([]byte(`{"age":0e-8, "created":100000e-4, "email":"foo@example.com", "id":0E-0, "name":"foo"}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		t.Parallel()

		var v user

		if err := Unmarshal(data, &v); err != nil {
			t.Skip()
		} else {
			t.Logf("Success decoding: %s", string(data))
		}
	})
}

func FuzzUnmarshalFields(f *testing.F) {
	f.Add(1, uint64(1), float64(1), "foo", "foo@example.com")
	f.Fuzz(func(t *testing.T, id int, created uint64, age float64, name, email string) {
		t.Parallel()

		data := marshalHelper(t, id, created, age, name, email)

		if r := recover(); r != nil {
			t.Fatalf("Unmarshal failed to properly handle: %s", data)
		}

		var v user

		if err := Unmarshal(data, &v); err != nil {
			t.Skip()
		} else {
			t.Logf("Success decoding: %s", string(data))
		}
	})
}

func marshalHelper(t *testing.T, id int, created uint64, age float64, name string, email string) []byte {
	t.Helper()

	data, err := json.Marshal(&user{ID: id, Created: created, Age: age, Name: name, Email: email})
	if err != nil {
		t.Fatalf("Unmarshal failed to properly handle: %d, %d, %f, %s, %s", id, created, age, name, email)
	}

	return data
}
