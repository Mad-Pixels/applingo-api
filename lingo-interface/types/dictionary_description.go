package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/pkg/errors"
)

// DictionaryDescription represents predefined dictionary descriptions
type DictionaryDescription int

const (
	DescBasicVocab DictionaryDescription = iota
	DescEverydayPhrases
	DescCommonExpressions
	DescEssentialWords
	DescPracticalVocab
	DescThematicCollection
	DescPopularTerms
	DescKeyPhrases
	DescUsefulExpressions
	DescCoreConcepts
)

// String returns the description text
func (d DictionaryDescription) String() string {
	switch d {
	case DescBasicVocab:
		return "Essential vocabulary for daily communication and basic conversations"
	case DescEverydayPhrases:
		return "Common phrases and expressions used in everyday situations"
	case DescCommonExpressions:
		return "Frequently used expressions and idioms for natural communication"
	case DescEssentialWords:
		return "Key words and terms necessary for effective language use"
	case DescPracticalVocab:
		return "Practical vocabulary for real-world situations and contexts"
	case DescThematicCollection:
		return "Thematic collection of words and phrases for specific topics"
	case DescPopularTerms:
		return "Popular and widely used terms in modern communication"
	case DescKeyPhrases:
		return "Key phrases and expressions for confident language usage"
	case DescUsefulExpressions:
		return "Useful expressions and vocabulary for various situations"
	case DescCoreConcepts:
		return "Core concepts and terminology for comprehensive understanding"
	default:
		return "General vocabulary collection for language learning"
	}
}

// AllDictionaryDescriptions returns a slice of all available descriptions
func AllDictionaryDescriptions() []DictionaryDescription {
	return []DictionaryDescription{
		DescBasicVocab,
		DescEverydayPhrases,
		DescCommonExpressions,
		DescEssentialWords,
		DescPracticalVocab,
		DescThematicCollection,
		DescPopularTerms,
		DescKeyPhrases,
		DescUsefulExpressions,
		DescCoreConcepts,
	}
}

// GetRandomDictionaryDescription returns a random dictionary description
func GetRandomDictionaryDescription() (DictionaryDescription, error) {
	descriptions := AllDictionaryDescriptions()
	idx, err := utils.RandomInt(0, len(descriptions)-1)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get random description")
	}
	return descriptions[idx], nil
}

// ParseDictionaryDescription converts string to DictionaryDescription
func ParseDictionaryDescription(s string) (DictionaryDescription, error) {
	for _, desc := range AllDictionaryDescriptions() {
		if desc.String() == s {
			return desc, nil
		}
	}
	return 0, errors.New("invalid dictionary description")
}
