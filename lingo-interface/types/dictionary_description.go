package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/pkg/errors"
)

// DictionaryDescription represents predefined dictionary descriptions
type DictionaryDescription int

const (
	DescEssentialCollection DictionaryDescription = iota
	DescPracticalTerms
	DescComprehensiveGuide
	DescLanguageToolkit
	DescVocabularyEssentials
	DescExpressionIndex
	DescLanguageResource
	DescWordCollection
	DescLinguisticCompendium
	DescCommunicationEssentials
	DescPhraseRepository
	DescVocabularyCompilation
	DescLearningCompanion
	DescExpressionCatalog
	DescWordInventory
	DescLanguageElements
	DescPhraseDirectory
	DescTerminologyGuide
	DescVocabularySpectrum
	DescExpressiveToolbox
	DescLanguageFoundation
	DescWordPortfolio
	DescConversationalTools
	DescLinguisticSelection
	DescVocabularyPanorama
)

// String returns the description text
func (d DictionaryDescription) String() string {
	switch d {
	case DescEssentialCollection:
		return "A carefully curated collection of essential words and expressions for effective communication"
	case DescPracticalTerms:
		return "Practical terminology and phrases for real-world conversations and situations"
	case DescComprehensiveGuide:
		return "A comprehensive guide to important vocabulary and meaningful expressions"
	case DescLanguageToolkit:
		return "A versatile toolkit of words and phrases for diverse communication needs"
	case DescVocabularyEssentials:
		return "Essential vocabulary selections to enhance language fluency and understanding"
	case DescExpressionIndex:
		return "An index of useful expressions and terminology for natural communication"
	case DescLanguageResource:
		return "A valuable resource of words and phrases for language development"
	case DescWordCollection:
		return "A thoughtfully assembled collection of words and expressions for daily use"
	case DescLinguisticCompendium:
		return "A compendium of linguistic elements to enrich your language skills"
	case DescCommunicationEssentials:
		return "Communication essentials for meaningful interactions and conversations"
	case DescPhraseRepository:
		return "A repository of phrases and vocabulary for effective self-expression"
	case DescVocabularyCompilation:
		return "A compilation of vocabulary designed to enhance communication abilities"
	case DescLearningCompanion:
		return "A companion of words and expressions to support language learning journey"
	case DescExpressionCatalog:
		return "A catalog of expressions and terminology for authentic communication"
	case DescWordInventory:
		return "An inventory of words and phrases to build confidence in language use"
	case DescLanguageElements:
		return "Fundamental language elements for clear and precise communication"
	case DescPhraseDirectory:
		return "A directory of phrases and terminology for expanding language capabilities"
	case DescTerminologyGuide:
		return "A guide to terminology and expressions for versatile language application"
	case DescVocabularySpectrum:
		return "A spectrum of vocabulary to enhance expression and comprehension"
	case DescExpressiveToolbox:
		return "A toolbox of expressive words and phrases for diverse communication contexts"
	case DescLanguageFoundation:
		return "Foundational language components for building strong communication skills"
	case DescWordPortfolio:
		return "A portfolio of words and expressions to enrich your language repertoire"
	case DescConversationalTools:
		return "Essential conversational tools for natural and fluid communication"
	case DescLinguisticSelection:
		return "A carefully selected linguistic collection for effective language use"
	case DescVocabularyPanorama:
		return "A panorama of vocabulary to broaden language horizons and abilities"
	default:
		return "A valuable collection of words and expressions for language learning"
	}
}

// AllDictionaryDescriptions returns a slice of all available descriptions
func AllDictionaryDescriptions() []DictionaryDescription {
	return []DictionaryDescription{
		DescEssentialCollection,
		DescPracticalTerms,
		DescComprehensiveGuide,
		DescLanguageToolkit,
		DescVocabularyEssentials,
		DescExpressionIndex,
		DescLanguageResource,
		DescWordCollection,
		DescLinguisticCompendium,
		DescCommunicationEssentials,
		DescPhraseRepository,
		DescVocabularyCompilation,
		DescLearningCompanion,
		DescExpressionCatalog,
		DescWordInventory,
		DescLanguageElements,
		DescPhraseDirectory,
		DescTerminologyGuide,
		DescVocabularySpectrum,
		DescExpressiveToolbox,
		DescLanguageFoundation,
		DescWordPortfolio,
		DescConversationalTools,
		DescLinguisticSelection,
		DescVocabularyPanorama,
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
