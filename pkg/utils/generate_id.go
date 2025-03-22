package utils

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	dictionaryFileExt = ".json"
)

// GenerateDictionaryID generates an MD5 hash from the concatenation of the dictionary name and author,
// separated by a hyphen. This hash serves as the unique object ID for an item in the DynamoDB dictionary table.
func GenerateDictionaryID(name, author string) string {
	return generateID(name, author)
}

// GenerateSubcategoryID generates an MD5 hash from the concatenation of the lang code and card side,
// separated by a hyphen.This hash serves as the unique object ID for an item in the DynamoDB subcategory table.
func GenerateSubcategoryID(code, side string) string {
	return generateID(code, side)
}

// RecordToFileID return file identifier based on incomming DynamoDB record ID.
func RecordToFileID(id string) string {
	switch {
	case id == "":
		return ""
	case filepath.Ext(id) != dictionaryFileExt:
		return id + dictionaryFileExt
	default:
		return id
	}
}

// IsFileID checks if the provided string is a valid file identifier
// A valid file ID must have a .json extension and the filename (without extension)
// must be a valid MD5 hash (32 hexadecimal characters)
func IsFileID(s string) bool {
	if s == "" {
		return true
	}
	if !strings.HasSuffix(s, dictionaryFileExt) {
		return false
	}

	prefix := strings.TrimSuffix(s, dictionaryFileExt)
	md5Regex := regexp.MustCompile(`^[a-f0-9]{32}$`)
	return md5Regex.MatchString(prefix)
}

func generateID(val1, val2 string) string {
	hash := md5.New()
	hash.Write([]byte(val1 + "-" + val2))
	return hex.EncodeToString(hash.Sum(nil))
}
