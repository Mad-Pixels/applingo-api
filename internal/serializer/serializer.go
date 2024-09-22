// Package serializer provides optimized JSON serialization and deserialization functions.
package serializer

import (
	"bytes"
	"encoding/json"
	"sync"
)

// bufPool is a pool of byte buffers used to reduce memory allocations.
var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

// MarshalJSON serializes the given value into a JSON-encoded byte slice.
// Example:
//
//	type Person struct {
//		Name string `json:"name"`
//		Age  int    `json:"age"`
//	}
//
//	p := Person{Name: "John Doe", Age: 30}
//	data, err := serializer.MarshalJSON(p)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(string(data))
//	// Output: {"name":"John Doe","age":30}
func MarshalJSON(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	// Remove the trailing newline added by enc.Encode
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// UnmarshalJSON deserializes the JSON-encoded data into the given value.
// Example:
//
//	data := []byte(`{"name":"Jane Doe","age":25}`)
//	var p Person
//	err := serializer.UnmarshalJSON(data, &p)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%+v\n", p)
//	// Output: {Name:Jane Doe Age:25}
func UnmarshalJSON(data []byte, v interface{}) error {
	buf := bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	buf.Write(data)
	return json.NewDecoder(buf).Decode(v)
}

// Bool is a custom boolean type that tracks if a value was set.
type Bool struct {
	Set   bool
	Value bool
}

// UnmarshalJSON implements the json.Unmarshaler interface for Bool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		b.Set = false
		return nil
	}
	err := json.Unmarshal(data, &b.Value)
	if err != nil {
		return err
	}
	b.Set = true
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Bool.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Set {
		return []byte("null"), nil
	}
	return json.Marshal(b.Value)
}
