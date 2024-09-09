package serializer

import (
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
			want:  &testStruct{Name: "John"},
		},
		{
			name:  "nested struct",
			input: `{"Person":{"Age":30}}`,
			want:  &nestedStruct{Person: struct{ Age int }{Age: 30}},
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
			switch tt.name {
			case "simple struct":
				got = &testStruct{}
			case "nested struct":
				got = &nestedStruct{}
			default:
				got = &map[string]interface{}{}
			}

			err := UnmarshalJSON([]byte(tt.input), got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
