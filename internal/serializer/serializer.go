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
		return new(bytes.Buffer)
	},
}

// MarshalJSON serializes the given value into a JSON-encoded byte slice.
// It uses a pool of buffers to reduce memory allocations and disables HTML escaping
// for better performance.
//
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
	return buf.Bytes(), nil
}

// UnmarshalJSON deserializes the JSON-encoded data into the given value.
// It uses a JSON decoder for efficient parsing of the input data.
//
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
	dec := json.NewDecoder(bytes.NewReader(data))
	return dec.Decode(v)
}