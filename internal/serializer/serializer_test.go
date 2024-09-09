package serializer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    string
		wantErr bool
	}{
		{
			name:  "simple struct",
			input: struct{ Name string }{Name: "John"},
			want:  `{"Name":"John"}`,
		},
		{
			name:  "nested struct",
			input: struct{ Person struct{ Age int } }{Person: struct{ Age int }{Age: 30}},
			want:  `{"Person":{"Age":30}}`,
		},
		{
			name:    "invalid input",
			input:   make(chan int),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalJSON(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, tt.want, string(got))
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type testStruct struct {
		Name string
	}

	type nestedStruct struct {
		Person struct {
			Age int
		}
	}

	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{
			name:  "simple struct",
			input: `{"Name":"John"}`,
			want:  testStruct{Name: "John"},
		},
		{
			name:  "nested struct",
			input: `{"Person":{"Age":30}}`,
			want:  nestedStruct{Person: struct{ Age int }{Age: 30}},
		},
		{
			name:    "invalid input",
			input:   `{"Name":"John"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			if tt.name == "simple struct" {
				got = &testStruct{}
			} else if tt.name == "nested struct" {
				got = &nestedStruct{}
			} else {
				got = &map[string]interface{}{}
			}

			err := UnmarshalJSON([]byte(tt.input), got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, *(got.(*struct{ Person struct{ Age int } })))
			}
		})
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	data := struct {
		Name string
		Age  int
		Tags []string
	}{
		Name: "John Doe",
		Age:  30,
		Tags: []string{"go", "programming", "json"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalJSON(data)
		require.NoError(b, err)
	}
}

func BenchmarkStandardMarshalJSON(b *testing.B) {
	data := struct {
		Name string
		Age  int
		Tags []string
	}{
		Name: "John Doe",
		Age:  30,
		Tags: []string{"go", "programming", "json"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		require.NoError(b, err)
	}
}