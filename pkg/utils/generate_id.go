package utils

import (
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
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
	case filepath.Ext(id) != ".json":
		return id + ".json"
	default:
		return id
	}
}

func generateID(val1, val2 string) string {
	hash := md5.New()
	hash.Write([]byte(val1 + "-" + val2))
	return hex.EncodeToString(hash.Sum(nil))
}
