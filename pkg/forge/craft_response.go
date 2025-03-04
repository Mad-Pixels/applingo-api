package forge

// DictionaryMetaFromAI represents metadata about the dictionary generated by the AI.
// It includes a brief description, the author responsible for the content, and the dictionary's name.
type DictionaryMetaFromAI struct {
	// Description provides an overview or summary of the generated dictionary.
	Description string `json:"description"`
	// Author indicates the creator or the AI that generated the dictionary.
	Author string `json:"author"`
	// Name specifies the title or name of the generated dictionary.
	Name string `json:"name"`
}

// DictionaryWordFromAI represents a single word entry within the generated dictionary.
// It contains the word itself, its translation, an explanation or description of the word, and an optional hint.
type DictionaryWordFromAI struct {
	// Description offers additional context or meaning for the word entry.
	Description string `json:"description"`
	// Translation provides the translated version of the word.
	Translation string `json:"translation"`
	// Word represents the vocabulary entry or term in the dictionary.
	Word string `json:"word"`
	// Hint gives an optional clue or additional tip related to the word.
	Hint string `json:"hint"`
}

// ResponseDictionaryCraft represents the complete response payload after generating a dictionary.
// It includes the dictionary metadata, a list of word entries, and a reference to the original request parameters.
type ResponseDictionaryCraft struct {
	// Request references the original dictionary craft request parameters that were used to generate the dictionary.
	Data *DictionaryCraftData `json:"data"`
	// Meta holds the metadata information such as description, author, and name of the generated dictionary.
	Meta DictionaryMetaFromAI `json:"meta"`
	// Words is a slice of dictionary word entries generated by the AI.
	Words []DictionaryWordFromAI `json:"words"`
}
