package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// PersistentID represents a unique, reproducible, and decomposable identifier.
// It consists of three parts: a timestamp, a unique ID, and a checksum.
type PersistentID struct {
	Timestamp int64
	UniqueID  string
	Checksum  string
}

// NewPersistentID creates a new PersistentID using the provided unique data.
// It generates a timestamp, creates a unique ID based on the input data and timestamp,
// and calculates a checksum for integrity verification.
//
// Example usage:
//
//	uniqueData := "user123@example.com"
//	pid := NewPersistentID(uniqueData)
//	fmt.Println(pid.String()) // Output: 1631234567890123456-a1b2c3d4e5f6g7h8-1234
func NewPersistentID(uniqueData string) PersistentID {
	timestamp := time.Now().UnixNano()
	uniqueID := generateUniqueID(uniqueData, timestamp)
	checksum := calculateChecksum(fmt.Sprintf("%d%s", timestamp, uniqueID))

	return PersistentID{
		Timestamp: timestamp,
		UniqueID:  uniqueID,
		Checksum:  checksum,
	}
}

// ParsePersistentID parses a string representation of a PersistentID and returns the corresponding struct.
// It expects the input string to be in the format: "timestamp-uniqueID-checksum".
//
// Example usage:
//
//	idString := "1631234567890123456-a1b2c3d4e5f6g7h8-1234"
//	pid, err := ParsePersistentID(idString)
//	if err != nil {
//	    fmt.Println("Error parsing PersistentID:", err)
//	    return
//	}
//	fmt.Printf("Timestamp: %d, UniqueID: %s, Checksum: %s\n", pid.Timestamp, pid.UniqueID, pid.Checksum)
func ParsePersistentID(idString string) (PersistentID, error) {
	var pid PersistentID
	_, err := fmt.Sscanf(idString, "%d-%s-%s", &pid.Timestamp, &pid.UniqueID, &pid.Checksum)
	if err != nil {
		return PersistentID{}, fmt.Errorf("invalid PersistentID format: %v", err)
	}
	return pid, nil
}

// String returns a string representation of the PersistentID.
// The format is: "timestamp-uniqueID-checksum".
func (pid PersistentID) String() string {
	return fmt.Sprintf("%d-%s-%s", pid.Timestamp, pid.UniqueID, pid.Checksum)
}

// Validate checks if the PersistentID is valid by recalculating the checksum
// and comparing it with the stored checksum.
//
// Example usage:
//
//	pid := NewPersistentID("user123@example.com")
//	if pid.Validate() {
//	    fmt.Println("PersistentID is valid")
//	} else {
//	    fmt.Println("PersistentID is invalid")
//	}
func (pid PersistentID) Validate() bool {
	expectedChecksum := calculateChecksum(fmt.Sprintf("%d%s", pid.Timestamp, pid.UniqueID))
	return pid.Checksum == expectedChecksum
}

func generateUniqueID(data string, timestamp int64) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", data, timestamp)))
	return hex.EncodeToString(hash[:16])
}

func calculateChecksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:4])
}
