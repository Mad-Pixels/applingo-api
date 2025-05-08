package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/pkg/errors"
)

// DictionaryTopic represents common dictionary topics.
type DictionaryTopic int

const (
	// TopicBasics represents basic vocabulary.
	TopicBasics DictionaryTopic = iota

	// TopicFood represents food and cooking.
	TopicFood

	// TopicTravel represents travel and tourism.
	TopicTravel

	// TopicBusiness represents business and management.
	TopicBusiness

	// TopicTechnology represents technology and innovation.
	TopicTechnology

	// TopicHealth represents health and medicine.
	TopicHealth

	// TopicSports represents sports and fitness.
	TopicSports

	// TopicEntertainment represents entertainment and media.
	TopicEntertainment

	// TopicScience represents science and research.
	TopicScience

	// TopicNature represents nature and environment.
	TopicNature

	// TopicArt represents art and creativity.
	TopicArt

	// TopicMusic represents music and performance.
	TopicMusic

	// TopicLiterature represents literature and writing.
	TopicLiterature

	// TopicMovies represents movies and cinematography.
	TopicMovies

	// TopicHistory represents history and civilizations.
	TopicHistory

	// TopicPsychology represents psychology and behavior.
	TopicPsychology

	// TopicPhilosophy represents philosophy and ethics.
	TopicPhilosophy

	// TopicEconomics represents economics and markets.
	TopicEconomics

	// TopicPolitics represents politics and governance.
	TopicPolitics

	// TopicGeography represents geography and cultures.
	TopicGeography

	// TopicAstronomy represents astronomy and space.
	TopicAstronomy

	// TopicPets represents pets and animal care.
	TopicPets

	// TopicGardening represents gardening and plants.
	TopicGardening

	// TopicArchitecture represents architecture and design.
	TopicArchitecture

	// TopicPhotography represents photography and visual arts.
	TopicPhotography

	// TopicRelationships represents relationships and communication.
	TopicRelationships

	// TopicInternet represents internet and digital life.
	TopicInternet

	// TopicEducation represents education and learning.
	TopicEducation

	// TopicFinance represents finance and investing.
	TopicFinance

	// TopicSocialMedia represents social media and networking.
	TopicSocialMedia

	// TopicInnovations represents innovations and discoveries.
	TopicInnovations

	// TopicLaw represents law and justice.
	TopicLaw

	// TopicMedicine represents medicine and healthcare.
	TopicMedicine

	// TopicFashion represents fashion and style.
	TopicFashion

	// TopicCulture represents culture and traditions.
	TopicCulture

	// TopicHobbies represents hobbies and leisure activities.
	TopicHobbies

	// TopicWorkplace represents workplace and professional life.
	TopicWorkplace

	// TopicHome represents home and household management.
	TopicHome

	// TopicTransportation represents transportation and travel modes.
	TopicTransportation

	// TopicWeather represents weather and climate.
	TopicWeather

	// TopicEmotions represents emotions and feelings.
	TopicEmotions

	// TopicHolidays represents holidays and celebrations.
	TopicHolidays

	// TopicShopping represents shopping and consumer behavior.
	TopicShopping

	// TopicAcademic represents academic and research terminology.
	TopicAcademic

	// TopicUrbanLife represents urban life and city navigation.
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
