package types

import (
	"strings"

	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/pkg/errors"
)

// LanguageLevel represents CEFR language proficiency level
type LanguageLevel int

const (
	// Beginner levels
	LevelA1 LanguageLevel = iota // Basic
	LevelA2                      // Elementary

	// Intermediate levels
	LevelB1 // Intermediate
	LevelB2 // Upper Intermediate

	// Advanced levels
	LevelC1 // Advanced
	LevelC2 // Mastery
)

// String returns the string representation of the language level
func (l LanguageLevel) String() string {
	switch l {
	case LevelA1:
		return "A1"
	case LevelA2:
		return "A2"
	case LevelB1:
		return "B1"
	case LevelB2:
		return "B2"
	case LevelC1:
		return "C1"
	case LevelC2:
		return "C2"
	default:
		return "undefined"
	}
}

// AllLanguageLevels returns a slice of all available language levels
func AllLanguageLevels() []LanguageLevel {
	return []LanguageLevel{
		LevelA1,
		LevelA2,
		LevelB1,
		LevelB2,
		LevelC1,
		LevelC2,
	}
}

// GetRandomLanguageLevel returns a random language level
func GetRandomLanguageLevel() (LanguageLevel, error) {
	levels := AllLanguageLevels()
	idx, err := utils.RandomInt(0, len(levels)-1)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get random language level")
	}
	return levels[idx], nil
}

// ParseLanguageLevel converts string to LanguageLevel
func ParseLanguageLevel(s string) (LanguageLevel, error) {
	s = strings.ToUpper(s)
	switch s {
	case "A1":
		return LevelA1, nil
	case "A2":
		return LevelA2, nil
	case "B1":
		return LevelB1, nil
	case "B2":
		return LevelB2, nil
	case "C1":
		return LevelC1, nil
	case "C2":
		return LevelC2, nil
	default:
		return 0, errors.New("invalid language level")
	}
}
