package types

import (
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/pkg/errors"
)

// DictionaryTopic represents practical, real-life dictionary topics.
type DictionaryTopic int

const (
	// ----------------------------- >
	// Group: Travel & Transportation
	// ----------------------------- >

	// TopicAirportTravel represents vocabulary related to airports and air travel.
	TopicAirportTravel DictionaryTopic = iota

	// TopicHotel represents vocabulary related to hotels and accommodations.
	TopicHotel

	// TopicDiningOut represents vocabulary for restaurants and dining out.
	TopicDiningOut

	// TopicGroceryShopping represents vocabulary for supermarkets and grocery stores.
	TopicGroceryShopping

	// TopicPublicTransport represents vocabulary for public transportation and tickets.
	TopicPublicTransport

	// TopicNavigation represents vocabulary for navigating streets and asking directions.
	TopicNavigation

	// TopicCarRental represents vocabulary for renting and using cars abroad.
	TopicCarRental

	// TopicSightseeing represents vocabulary for tourist attractions and sightseeing.
	TopicSightseeing

	// TopicTransportApps represents vocabulary for using local transport apps and maps.
	TopicTransportApps

	// TopicCityNavigation represents vocabulary for navigating urban environments.
	TopicCityNavigation

	// Group: Emergencies & Healthcare

	// TopicEmergency represents vocabulary for emergencies and healthcare situations.
	TopicEmergency

	// TopicPharmacy represents vocabulary for pharmacies and medications.
	TopicPharmacy

	// TopicMedicalVisit represents vocabulary for visiting doctors and describing symptoms.
	TopicMedicalVisit

	// TopicVetVisit represents vocabulary for pets and veterinary visits.
	TopicVetVisit

	// ----------------------------- >
	// Group: Daily Life & Services
	// ----------------------------- >

	// TopicCafeCulture represents vocabulary related to cafés and coffee culture.
	TopicCafeCulture

	// TopicLocalFoods represents vocabulary for local food items and bakeries.
	TopicLocalFoods

	// TopicSeasonalClothing represents vocabulary for clothing and weather conditions.
	TopicSeasonalClothing

	// TopicSIMCards represents vocabulary for mobile phones and SIM cards abroad.
	TopicSIMCards

	// TopicRealEstate represents vocabulary for renting and buying property.
	TopicRealEstate

	// TopicSpecialDiets represents vocabulary for food allergies and dietary preferences.
	TopicSpecialDiets

	// TopicFestivals represents vocabulary related to festivals and local holidays.
	TopicFestivals

	// TopicPublicServices represents vocabulary for city services and public offices.
	TopicPublicServices

	// TopicStreetFood represents vocabulary for local markets and street food.
	TopicStreetFood

	// TopicHairSalon represents vocabulary for visiting barbershops and hair salons.
	TopicHairSalon

	// TopicWeather represents vocabulary for weather forecasts and natural conditions.
	TopicWeather

	// TopicPostOffice represents vocabulary for postal services and parcel delivery.
	TopicPostOffice

	// TopicLibrary represents vocabulary for libraries and studying.
	TopicLibrary

	// TopicRecycling represents vocabulary for waste sorting and recycling.
	TopicRecycling

	// TopicLostAndFound represents vocabulary for lost items and related conversations.
	TopicLostAndFound

	// ----------------------------- >
	// Group: Work & Business
	// ----------------------------- >

	// TopicJobInterviews represents vocabulary for job interviews and resumes.
	TopicJobInterviews

	// TopicOfficeCommunication represents vocabulary for communication in office settings.
	TopicOfficeCommunication

	// TopicBusinessWriting represents vocabulary for emails and formal writing.
	TopicBusinessWriting

	// TopicRemoteWork represents vocabulary for remote work and freelancing.
	TopicRemoteWork

	// TopicOnlineShopping represents vocabulary for e-commerce and deliveries.
	TopicOnlineShopping

	// TopicBanking represents vocabulary for banking and opening accounts.
	TopicBanking

	// TopicInvesting represents vocabulary for personal finance and investing.
	TopicInvesting

	// TopicInsurance represents vocabulary for insurance policies and claims.
	TopicInsurance

	// TopicDigitalSecurity represents vocabulary related to cybersecurity and privacy.
	TopicDigitalSecurity

	// ----------------------------- >
	// Group: Communication & Culture
	// ----------------------------- >

	// TopicCulturalEtiquette represents vocabulary for customs, etiquette, and manners.
	TopicCulturalEtiquette

	// TopicSocialMedia represents vocabulary for social networks and digital profiles.
	TopicSocialMedia

	// TopicConnectivity represents vocabulary for public Wi-Fi and internet access.
	TopicConnectivity

	// TopicSmallTalk represents vocabulary for small talk and casual conversations.
	TopicSmallTalk

	// TopicLanguageExchange represents vocabulary for language exchange and studying abroad.
	TopicLanguageExchange

	// ----------------------------- >
	// Group: Hobbies & Leisure
	// ----------------------------- >

	// TopicMuseums represents vocabulary for museums and art galleries.
	TopicMuseums

	// TopicTravelPhotography represents vocabulary for taking photos while traveling.
	TopicTravelPhotography

	// TopicGardening represents vocabulary for gardening tools and plant care.
	TopicGardening

	// TopicDIY represents vocabulary for home repairs and do-it-yourself projects.
	TopicDIY

	// TopicHolidays represents vocabulary for holiday activities and celebrations.
	TopicHolidays

	// TopicShopping represents vocabulary for shopping and consumer behavior.
	TopicShopping
)

// String returns the string representation of the topic.
func (t DictionaryTopic) String() string {
	switch t {
	case TopicAirportTravel:
		return "Airport and Air Travel"
	case TopicHotel:
		return "Hotel and Accommodation"
	case TopicDiningOut:
		return "Restaurants and Dining Out"
	case TopicGroceryShopping:
		return "Supermarket and Grocery Shopping"
	case TopicEmergency:
		return "Emergency and Healthcare Situations"
	case TopicPublicTransport:
		return "Public Transportation and Tickets"
	case TopicNavigation:
		return "Street Navigation and Directions"
	case TopicCarRental:
		return "Renting a Car Abroad"
	case TopicJobInterviews:
		return "Job Interviews and Resumes"
	case TopicOfficeCommunication:
		return "Office Communication and Meetings"
	case TopicBusinessWriting:
		return "Email and Business Writing"
	case TopicOnlineShopping:
		return "Online Shopping and Delivery"
	case TopicDigitalSecurity:
		return "Digital Security and Privacy"
	case TopicBanking:
		return "Banking and Opening an Account"
	case TopicInvesting:
		return "Investing and Personal Finance"
	case TopicInsurance:
		return "Insurance and Claims"
	case TopicPharmacy:
		return "Pharmacy and Medication"
	case TopicMedicalVisit:
		return "Doctor’s Appointment and Symptoms"
	case TopicCafeCulture:
		return "Cafés and Coffee Culture"
	case TopicLocalFoods:
		return "Bakeries and Local Foods"
	case TopicSightseeing:
		return "Tourist Attractions and Sightseeing"
	case TopicSeasonalClothing:
		return "Climate and Seasonal Clothing"
	case TopicTransportApps:
		return "Local Transportation Apps and Maps"
	case TopicCulturalEtiquette:
		return "Cultural Etiquette and Manners"
	case TopicSIMCards:
		return "Mobile Phones and SIM Cards Abroad"
	case TopicRealEstate:
		return "Real Estate and Apartment Search"
	case TopicSpecialDiets:
		return "Food Allergies and Special Diets"
	case TopicFestivals:
		return "Festivals and Local Holidays"
	case TopicPublicServices:
		return "City Services and Public Offices"
	case TopicRemoteWork:
		return "Working Remotely and Freelancing"
	case TopicMuseums:
		return "Art Galleries and Museums"
	case TopicTravelPhotography:
		return "Photography and Taking Pictures Abroad"
	case TopicSocialMedia:
		return "Social Media and Online Presence"
	case TopicConnectivity:
		return "Public Wi-Fi and Connectivity"
	case TopicSmallTalk:
		return "Small Talk and Daily Conversations"
	case TopicVetVisit:
		return "Pets and Vet Visits"
	case TopicGardening:
		return "Gardening Tools and Plants"
	case TopicStreetFood:
		return "Local Markets and Street Food"
	case TopicDIY:
		return "DIY and Home Repairs"
	case TopicHairSalon:
		return "Hairdresser and Barbershop Vocabulary"
	case TopicWeather:
		return "Weather Forecast and Natural Conditions"
	case TopicLibrary:
		return "Libraries and Study Materials"
	case TopicPostOffice:
		return "Parcel Pickup and Postal Services"
	case TopicLanguageExchange:
		return "Language Exchange and Study Abroad"
	case TopicRecycling:
		return "Recycling and Waste Sorting"
	case TopicLostAndFound:
		return "Lost and Found Situations"
	case TopicHolidays:
		return "Holidays and Celebrations"
	case TopicShopping:
		return "Shopping and Consumer Behavior"
	case TopicCityNavigation:
		return "Urban Life and City Navigation"
	default:
		return "Unknown Topic"
	}
}

// AllDictionaryTopics returns a slice of all available dictionary topics.
func AllDictionaryTopics() []DictionaryTopic {
	return []DictionaryTopic{
		TopicAirportTravel, TopicHotel, TopicDiningOut, TopicGroceryShopping,
		TopicEmergency, TopicPublicTransport, TopicNavigation, TopicCarRental,
		TopicJobInterviews, TopicOfficeCommunication, TopicBusinessWriting,
		TopicOnlineShopping, TopicDigitalSecurity, TopicBanking, TopicInvesting,
		TopicInsurance, TopicPharmacy, TopicMedicalVisit,
		TopicCafeCulture, TopicLocalFoods, TopicSightseeing, TopicSeasonalClothing,
		TopicTransportApps, TopicCulturalEtiquette, TopicSIMCards, TopicRealEstate,
		TopicSpecialDiets, TopicFestivals, TopicPublicServices, TopicRemoteWork,
		TopicMuseums, TopicTravelPhotography, TopicSocialMedia, TopicConnectivity,
		TopicSmallTalk, TopicVetVisit, TopicGardening, TopicStreetFood,
		TopicDIY, TopicHairSalon, TopicWeather, TopicLibrary, TopicPostOffice,
		TopicLanguageExchange, TopicRecycling, TopicLostAndFound, TopicHolidays,
		TopicShopping, TopicCityNavigation,
	}
}

// GetRandomDictionaryTopic returns a random dictionary topic.
func GetRandomDictionaryTopic() (DictionaryTopic, error) {
	topics := AllDictionaryTopics()
	idx, err := utils.RandomInt(0, len(topics)-1)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get random topic")
	}
	return topics[idx], nil
}

// ParseDictionaryTopic converts a string to DictionaryTopic.
func ParseDictionaryTopic(s string) (DictionaryTopic, error) {
	for _, topic := range AllDictionaryTopics() {
		if topic.String() == s {
			return topic, nil
		}
	}
	return 0, errors.New("invalid dictionary topic")
}
