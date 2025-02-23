package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/pkg/errors"
)

// DictionaryTopic represents common dictionary topics
type DictionaryTopic int

const (
	TopicBasics DictionaryTopic = iota
	TopicFood
	TopicTravel
	TopicBusiness
	TopicTechnology
	TopicEducation
	TopicHealth
	TopicSports
	TopicEntertainment
	TopicNature
	TopicScience
	TopicArt
	TopicMusic
	TopicFamily
	TopicWork
	TopicHobbies
	TopicClothes
	TopicEmotions
	TopicWeather
	TopicTransport
)

// String returns the string representation of the topic
func (t DictionaryTopic) String() string {
	switch t {
	case TopicBasics:
		return "Basic Vocabulary"
	case TopicFood:
		return "Food and Cooking"
	case TopicTravel:
		return "Travel and Places"
	case TopicBusiness:
		return "Business and Finance"
	case TopicTechnology:
		return "Technology and Internet"
	case TopicEducation:
		return "Education and Learning"
	case TopicHealth:
		return "Health and Medicine"
	case TopicSports:
		return "Sports and Fitness"
	case TopicEntertainment:
		return "Entertainment and Media"
	case TopicNature:
		return "Nature and Environment"
	case TopicScience:
		return "Science and Research"
	case TopicArt:
		return "Art and Culture"
	case TopicMusic:
		return "Music and Performance"
	case TopicFamily:
		return "Family and Relationships"
	case TopicWork:
		return "Work and Career"
	case TopicHobbies:
		return "Hobbies and Activities"
	case TopicClothes:
		return "Clothes and Fashion"
	case TopicEmotions:
		return "Emotions and Feelings"
	case TopicWeather:
		return "Weather and Climate"
	case TopicTransport:
		return "Transport and Travel"
	default:
		return "Unknown Topic"
	}
}

// AllDictionaryTopics returns a slice of all available topics
func AllDictionaryTopics() []DictionaryTopic {
	return []DictionaryTopic{
		TopicBasics, TopicFood, TopicTravel, TopicBusiness,
		TopicTechnology, TopicEducation, TopicHealth, TopicSports,
		TopicEntertainment, TopicNature, TopicScience, TopicArt,
		TopicMusic, TopicFamily, TopicWork, TopicHobbies,
		TopicClothes, TopicEmotions, TopicWeather, TopicTransport,
	}
}

// GetRandomDictionaryTopic returns a random dictionary topic
func GetRandomDictionaryTopic() (DictionaryTopic, error) {
	topics := AllDictionaryTopics()
	idx, err := utils.RandomInt(0, len(topics)-1)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get random topic")
	}
	return topics[idx], nil
}

// ParseDictionaryTopic converts string to DictionaryTopic
func ParseDictionaryTopic(s string) (DictionaryTopic, error) {
	for _, topic := range AllDictionaryTopics() {
		if topic.String() == s {
			return topic, nil
		}
	}
	return 0, errors.New("invalid dictionary topic")
}
