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
	TopicHealth
	TopicSports
	TopicEntertainment
	TopicScience
	TopicNature
	TopicArt
	TopicMusic
	TopicLiterature
	TopicMovies
	TopicHistory
	TopicPsychology
	TopicPhilosophy
	TopicEconomics
	TopicPolitics
	TopicGeography
	TopicAstronomy
	TopicPets
	TopicGardening
	TopicArchitecture
	TopicPhotography
	TopicRelationships
	TopicInternet
	TopicEducation
	TopicFinance
	TopicSocialMedia
	TopicInnovations
	TopicLaw
	TopicMedicine
	TopicFashion
	TopicCulture
	TopicHobbies
	TopicWorkplace
	TopicHome
	TopicTransportation
	TopicWeather
	TopicEmotions
	TopicHolidays
	TopicShopping
	TopicAcademic
	TopicUrbanLife
)

// String returns the string representation of the topic
func (t DictionaryTopic) String() string {
	switch t {
	case TopicBasics:
		return "Basic Vocabulary"
	case TopicFood:
		return "Food and Cooking"
	case TopicTravel:
		return "Travel and Tourism"
	case TopicBusiness:
		return "Business and Management"
	case TopicTechnology:
		return "Technology and Innovation"
	case TopicHealth:
		return "Health and Medicine"
	case TopicSports:
		return "Sports and Fitness"
	case TopicEntertainment:
		return "Entertainment and Media"
	case TopicScience:
		return "Science and Research"
	case TopicNature:
		return "Nature and Environment"
	case TopicArt:
		return "Art and Creativity"
	case TopicMusic:
		return "Music and Performance"
	case TopicLiterature:
		return "Literature and Writing"
	case TopicMovies:
		return "Movies and Cinematography"
	case TopicHistory:
		return "History and Civilizations"
	case TopicPsychology:
		return "Psychology and Behavior"
	case TopicPhilosophy:
		return "Philosophy and Ethics"
	case TopicEconomics:
		return "Economics and Markets"
	case TopicPolitics:
		return "Politics and Governance"
	case TopicGeography:
		return "Geography and Cultures"
	case TopicAstronomy:
		return "Astronomy and Space"
	case TopicPets:
		return "Pets and Animal Care"
	case TopicGardening:
		return "Gardening and Plants"
	case TopicArchitecture:
		return "Architecture and Design"
	case TopicPhotography:
		return "Photography and Visual Arts"
	case TopicRelationships:
		return "Relationships and Communication"
	case TopicInternet:
		return "Internet and Digital Life"
	case TopicEducation:
		return "Education and Learning"
	case TopicFinance:
		return "Finance and Investing"
	case TopicSocialMedia:
		return "Social Media and Networking"
	case TopicInnovations:
		return "Innovations and Discoveries"
	case TopicLaw:
		return "Law and Justice"
	case TopicMedicine:
		return "Medicine and Healthcare"
	case TopicFashion:
		return "Fashion and Style"
	case TopicCulture:
		return "Culture and Traditions"
	case TopicHobbies:
		return "Hobbies and Leisure Activities"
	case TopicWorkplace:
		return "Workplace and Professional Life"
	case TopicHome:
		return "Home and Household Management"
	case TopicTransportation:
		return "Transportation and Travel Modes"
	case TopicWeather:
		return "Weather and Climate"
	case TopicEmotions:
		return "Emotions and Feelings"
	case TopicHolidays:
		return "Holidays and Celebrations"
	case TopicShopping:
		return "Shopping and Consumer Behavior"
	case TopicAcademic:
		return "Academic and Research Terminology"
	case TopicUrbanLife:
		return "Urban Life and City Navigation"
	default:
		return "Unknown Topic"
	}
}

// AllDictionaryTopics returns a slice of all available topics
func AllDictionaryTopics() []DictionaryTopic {
	return []DictionaryTopic{
		TopicBasics, TopicFood, TopicTravel, TopicBusiness, TopicTechnology,
		TopicHealth, TopicSports, TopicEntertainment, TopicScience,
		TopicNature, TopicArt, TopicMusic, TopicLiterature, TopicMovies,
		TopicHistory, TopicPsychology, TopicPhilosophy, TopicEconomics,
		TopicPolitics, TopicGeography, TopicAstronomy, TopicPets,
		TopicGardening, TopicArchitecture, TopicPhotography, TopicRelationships,
		TopicInternet, TopicEducation, TopicFinance, TopicSocialMedia,
		TopicInnovations, TopicLaw, TopicMedicine, TopicFashion,
		TopicCulture, TopicHobbies, TopicWorkplace, TopicHome,
		TopicTransportation, TopicWeather, TopicEmotions, TopicHolidays,
		TopicShopping, TopicAcademic, TopicUrbanLife,
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
