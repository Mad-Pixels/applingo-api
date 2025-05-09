package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/pkg/errors"
)

// DictionaryDescription represents predefined dictionary descriptions.
type DictionaryDescription int

const (
	// ----------------------------- >
	// Group: Universal Descriptions
	// ----------------------------- >

	// DescPracticalUsage emphasizes practical usage and clarity.
	DescPracticalUsage DictionaryDescription = iota

	// DescEssentialPatterns focuses on essential communication patterns.
	DescEssentialPatterns

	// DescContextualExamples provides varied examples and contextual clarity.
	DescContextualExamples

	// DescRetentionFriendly structured for easier retention and recall.
	DescRetentionFriendly

	// DescRealWorldUsage organized to reflect real-world usage.
	DescRealWorldUsage

	// DescActivePassive designed to improve both active and passive vocabulary.
	DescActivePassive

	// DescFunctionalFocus highlights functionally important words.
	DescFunctionalFocus

	// DescContextUnderstanding promotes better understanding through context.
	DescContextUnderstanding

	// DescRegisterBalance balanced between formal and informal registers.
	DescRegisterBalance

	// DescDailyComprehension optimized for everyday comprehension and use.
	DescDailyComprehension

	// DescClearDefinitions with clear definitions and everyday contexts.
	DescClearDefinitions

	// DescConsistentAcquisition intended to support consistent vocabulary acquisition.
	DescConsistentAcquisition

	// DescFrequencyClarity prioritizes frequency and clarity over complexity.
	DescFrequencyClarity

	// DescMeaningfulAssociations curated for meaningful word associations.
	DescMeaningfulAssociations

	// DescContextLearning organized to encourage contextual learning.
	DescContextLearning

	// DescMemorabilityFocus built around usage relevance and memorability.
	DescMemorabilityFocus

	// DescFamiliarExamples encouraging word usage through familiar examples.
	DescFamiliarExamples

	// DescNaturalAbsorption optimized for natural language absorption.
	DescNaturalAbsorption

	// DescRetentionBoost aimed at reinforcing thematic vocabulary retention.
	DescRetentionBoost

	// DescGeneralCommunication tailored to general real-life communication.
	DescGeneralCommunication

	// DescCommunicativeValue featuring language with high communicative value.
	DescCommunicativeValue

	// DescWellRounded understanding of each word.
	DescWellRounded

	// DescRepetitionVariation strengthen comprehension through repetition and variation.
	DescRepetitionVariation

	// DescRelatableContexts framed around clear and relatable contexts.
	DescRelatableContexts

	// DescUniversalApplicability designed for universal applicability across situations.
	DescUniversalApplicability

	// DescEngagementRetention intended to maximize engagement and memory retention.
	DescEngagementRetention

	// DescFunctionBalance with focus on word function and usage balance.
	DescFunctionBalance

	// DescFluencyDevelopment supporting development of fluency through exposure.
	DescFluencyDevelopment

	// DescUsageScenarios with attention to common usage scenarios.
	DescUsageScenarios

	// DescReliableReference curated to serve as a reliable vocabulary reference.
	DescReliableReference
)

// String returns the description text.
func (d DictionaryDescription) String() string {
	switch d {
	case DescPracticalUsage:
		return "Emphasizing practical usage and clarity"
	case DescEssentialPatterns:
		return "Focusing on essential communication patterns"
	case DescContextualExamples:
		return "Providing varied examples and contextual clarity"
	case DescRetentionFriendly:
		return "Structured for easier retention and recall"
	case DescRealWorldUsage:
		return "Organized to reflect real-world usage"
	case DescActivePassive:
		return "Designed to improve both active and passive vocabulary"
	case DescFunctionalFocus:
		return "Highlighting functionally important words"
	case DescContextUnderstanding:
		return "Promoting better understanding through context"
	case DescRegisterBalance:
		return "Balanced between formal and informal registers"
	case DescDailyComprehension:
		return "Optimized for everyday comprehension and use"
	case DescClearDefinitions:
		return "With clear definitions and everyday contexts"
	case DescConsistentAcquisition:
		return "Intended to support consistent vocabulary acquisition"
	case DescFrequencyClarity:
		return "Prioritizing frequency and clarity over complexity"
	case DescMeaningfulAssociations:
		return "Curated for meaningful word associations"
	case DescContextLearning:
		return "Organized to encourage contextual learning"
	case DescMemorabilityFocus:
		return "Built around usage relevance and memorability"
	case DescFamiliarExamples:
		return "Encouraging word usage through familiar examples"
	case DescNaturalAbsorption:
		return "Optimized for natural language absorption"
	case DescRetentionBoost:
		return "Aimed at reinforcing thematic vocabulary retention"
	case DescGeneralCommunication:
		return "Tailored to general real-life communication"
	case DescCommunicativeValue:
		return "Featuring language with high communicative value"
	case DescWellRounded:
		return "Providing well-rounded understanding of each word"
	case DescRepetitionVariation:
		return "Built to strengthen comprehension through repetition and variation"
	case DescRelatableContexts:
		return "Framed around clear and relatable contexts"
	case DescUniversalApplicability:
		return "Designed for universal applicability across situations"
	case DescEngagementRetention:
		return "Intended to maximize engagement and memory retention"
	case DescFunctionBalance:
		return "With focus on word function and usage balance"
	case DescFluencyDevelopment:
		return "Supporting development of fluency through exposure"
	case DescUsageScenarios:
		return "With attention to common usage scenarios"
	case DescReliableReference:
		return "Curated to serve as a reliable vocabulary reference"
	default:
		return "General-purpose vocabulary description"
	}
}

// AllDictionaryDescriptions returns a slice of all available dictionary descriptions.
func AllDictionaryDescriptions() []DictionaryDescription {
	return []DictionaryDescription{
		DescPracticalUsage,
		DescEssentialPatterns,
		DescContextualExamples,
		DescRetentionFriendly,
		DescRealWorldUsage,
		DescActivePassive,
		DescFunctionalFocus,
		DescContextUnderstanding,
		DescRegisterBalance,
		DescDailyComprehension,
		DescClearDefinitions,
		DescConsistentAcquisition,
		DescFrequencyClarity,
		DescMeaningfulAssociations,
		DescContextLearning,
		DescMemorabilityFocus,
		DescFamiliarExamples,
		DescNaturalAbsorption,
		DescRetentionBoost,
		DescGeneralCommunication,
		DescCommunicativeValue,
		DescWellRounded,
		DescRepetitionVariation,
		DescRelatableContexts,
		DescUniversalApplicability,
		DescEngagementRetention,
		DescFunctionBalance,
		DescFluencyDevelopment,
		DescUsageScenarios,
		DescReliableReference,
	}
}

// GetRandomDictionaryDescription returns a random dictionary description.
func GetRandomDictionaryDescription() (DictionaryDescription, error) {
	descs := AllDictionaryDescriptions()
	idx, err := utils.RandomInt(0, len(descs)-1)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get random description")
	}
	return descs[idx], nil
}

// ParseDictionaryDescription converts a string to DictionaryDescription.
func ParseDictionaryDescription(s string) (DictionaryDescription, error) {
	for _, desc := range AllDictionaryDescriptions() {
		if desc.String() == s {
			return desc, nil
		}
	}
	return 0, errors.New("invalid dictionary description")
}
