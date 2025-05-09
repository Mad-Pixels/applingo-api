package types

import (
	"strings"

	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/pkg/errors"
)

// LanguageLevel represents a CEFR language proficiency level.
type LanguageLevel int

const (
	// LevelA1 represents the A1 (Beginner) language level.
	LevelA1 LanguageLevel = iota

	// LevelA2 represents the A2 (Elementary) language level.
	LevelA2

	// LevelB1 represents the B1 (Intermediate) language level.
	LevelB1

	// LevelB2 represents the B2 (Upper Intermediate) language level.
	LevelB2

	// LevelC1 represents the C1 (Advanced) language level.
	LevelC1

	// LevelC2 represents the C2 (Mastery) language level.
	LevelC2
)

// String returns the string representation of the LanguageLevel.
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

// AllLanguageLevels returns all supported CEFR language levels.
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

// GetRandomLanguageLevel returns a randomly selected CEFR language level.
func GetRandomLanguageLevel() (LanguageLevel, error) {
	levels := AllLanguageLevels()
	idx, err := utils.RandomInt(0, len(levels)-1)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get random language level")
	}
	return levels[idx], nil
}

// ParseLanguageLevel parses a string (e.g., "A2") into a LanguageLevel value.
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
