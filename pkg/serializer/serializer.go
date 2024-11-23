package serializer

import (
	"bytes"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            false,
	UseNumber:              false,
	ValidateJsonRawMessage: true,
}.Froze()

// bufPool is a pool of byte buffers used to reduce memory allocations.
var bufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024))
	},
}

// MarshalJSON serializes the given value into a JSON-encoded byte slice using jsoniter.
func MarshalJSON(v interface{}) ([]byte, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	enc := json.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// UnmarshalJSON deserializes the JSON-encoded data into the given value using jsoniter.
func UnmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
