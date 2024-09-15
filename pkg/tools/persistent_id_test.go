package tools

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPersistentID(t *testing.T) {
	tests := []struct {
		name       string
		uniqueData string
	}{
		{
			name:       "simple string",
			uniqueData: "test123",
		},
		{
			name:       "email address",
			uniqueData: "user@example.com",
		},
		{
			name:       "empty string",
			uniqueData: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid := NewPersistentID(tt.uniqueData)

			assert.NotEmpty(t, pid.Timestamp)
			assert.NotEmpty(t, pid.UniqueID)
			assert.NotEmpty(t, pid.Checksum)
			assert.True(t, pid.Validate())

			assert.WithinDuration(t, time.Now(), time.Unix(0, pid.Timestamp), 5*time.Second)
		})
	}
}

func TestPersistentIDString(t *testing.T) {
	pid := PersistentID{
		Timestamp: 1631234567890123456,
		UniqueID:  "a1b2c3d4e5f6g7h8",
		Checksum:  "1234",
	}

	expected := "1631234567890123456-a1b2c3d4e5f6g7h8-1234"
	assert.Equal(t, expected, pid.String())
}

func TestParsePersistentID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    PersistentID
		wantErr bool
	}{
		{
			name:  "valid PersistentID",
			input: "1631234567890123456-a1b2c3d4e5f6g7h8-1234",
			want: PersistentID{
				Timestamp: 1631234567890123456,
				UniqueID:  "a1b2c3d4e5f6g7h8",
				Checksum:  "1234",
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid-format",
			wantErr: true,
		},
		{
			name:    "missing parts",
			input:   "1631234567890123456-a1b2c3d4e5f6g7h8",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePersistentID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPersistentIDValidate(t *testing.T) {
	tests := []struct {
		name string
		pid  PersistentID
		want bool
	}{
		{
			name: "valid PersistentID",
			pid:  NewPersistentID("test123"),
			want: true,
		},
		{
			name: "invalid checksum",
			pid: PersistentID{
				Timestamp: time.Now().UnixNano(),
				UniqueID:  "a1b2c3d4e5f6g7h8",
				Checksum:  "invalid",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pid.Validate())
		})
	}
}

func TestGenerateUniqueID(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		timestamp int64
	}{
		{
			name:      "simple string",
			data:      "test123",
			timestamp: 1631234567890123456,
		},
		{
			name:      "empty string",
			data:      "",
			timestamp: 1631234567890123456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueID := generateUniqueID(tt.data, tt.timestamp)
			assert.Len(t, uniqueID, 32) // 16 bytes encoded as hex
			assert.Regexp(t, "^[0-9a-f]{32}$", uniqueID)
		})
	}
}

func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "simple string",
			data: "test123",
		},
		{
			name: "empty string",
			data: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checksum := calculateChecksum(tt.data)
			assert.Len(t, checksum, 8) // 4 bytes encoded as hex
			assert.Regexp(t, "^[0-9a-f]{8}$", checksum)
		})
	}
}
