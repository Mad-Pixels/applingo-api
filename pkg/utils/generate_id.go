package utils

import (
	"encoding/base64"
	"errors"
	"path/filepath"
	"strings"
)

// GenerateDictionaryID generates an MD5 hash from the concatenation of the dictionary name and author,
// separated by a hyphen. This hash serves as the unique object ID for an item in the DynamoDB dictionary table.
func GenerateDictionaryID(name, author string) string {
	return generateID(name, author)
}

// DecodeDictionaryID decodes an incoming ID to a []string containing the dictionary name and dictionary author.
func DecodeDictionaryID(id string) ([]string, error) {
	res, err := decodeID(id)
	if err != nil {
		return nil, errors.New("cannot decode id")
	}
	if len(res) != 2 {
		return nil, errors.New("decode result has invalid format")
	}
	return res, nil
}

// RecordToFileID return file identifier based on incomming DynamoDB record ID.
func RecordToFileID(id string) string {
	switch {
	case id == "":
		return ""
	case filepath.Ext(id) != ".json":
		return id + ".json"
	default:
		return id
	}
}

func generateID(val1, val2 string) string {
	data := val1 + "|" + val2
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func decodeID(id string) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(decoded), "|"), nil
}
