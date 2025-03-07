package forge

// CheckMetaFromAI represents metadata about the dictionary check generated by the AI.
// It includes a score and a reason.
type CheckMetaFromAI struct {
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

// ResponseDictionaryCheck represents the complete response payload for a dictionary check request.
// It includes the generated metadata and the original request parameters.
type ResponseDictionaryCheck struct {
	Meta CheckMetaFromAI      `json:"meta"`
	Data *DictionaryCheckData `json:"data"`
}
