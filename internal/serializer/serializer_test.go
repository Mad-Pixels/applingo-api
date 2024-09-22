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
		{
			name:  "custom Bool - set",
			input: struct{ Active Bool }{Active: Bool{Set: true, Value: true}},
			want:  `{"Active":true}`,
		},
		{
			name:  "custom Bool - not set",
			input: struct{ Active Bool }{Active: Bool{Set: false}},
			want:  `{"Active":null}`,
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

	type boolStruct struct {
		Active Bool
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
		{
			name:  "custom Bool - true",
			input: `{"Active":true}`,
			want:  &boolStruct{Active: Bool{Set: true, Value: true}},
		},
		{
			name:  "custom Bool - false",
			input: `{"Active":false}`,
			want:  &boolStruct{Active: Bool{Set: true, Value: false}},
		},
		{
			name:  "custom Bool - null",
			input: `{"Active":null}`,
			want:  &boolStruct{Active: Bool{Set: false}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got interface{}
			switch tt.want.(type) {
			case *testStruct:
				got = &testStruct{}
			case *nestedStruct:
				got = &nestedStruct{}
			case *boolStruct:
				got = &boolStruct{}
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
