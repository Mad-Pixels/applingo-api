package forge

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"

	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/google/uuid"
)

// craftDictionaryPromptTemplate defines the structure for the prompt used in dictionary generation.
// It aggregates the main parameters that will be injected into the AI prompt template.
type craftDictionaryPromptTemplate struct {
	// DictionaryDescription is a brief overview of the dictionary's content for the AI model.
	DictionaryDescription string
	// DictionaryTopic specifies the main subject or theme of the dictionary.
	DictionaryTopic string
	// LanguageLevel represents the CEFR proficiency level as a string (e.g., "A1", "B2").
	LanguageLevel string
	// LanguageFrom indicates the source language in string format.
	LanguageFrom string
	// LanguageTo indicates the target language in string format.
	LanguageTo string
	// WordsCount specifies the number of words to be generated in the dictionary.
	WordsCount int
}

// DictionaryCraftData holds all data related to a dictionary generation request.
// It encapsulates both the input parameters and internal data used during the crafting process.
type DictionaryCraftData struct {
	request   *RequestDictionaryCraft  // Original dictionary craft request parameters.
	response  *ResponseDictionaryCraft // Original openai craft response.
	promptBuf *bytes.Buffer            // Internal buffer to store the fetched prompt template data.

	prompt                string // Prompt string used for initiating dictionary generation.
	filename              string // Generated unique identifier (UUID) for the dictionary file.
	dictionaryDescription string // Common description of the dictionary for AI model.
	dictionaryTopic       string // Subject or theme of the dictionary for AI model.

	model         chatgpt.OpenAIModel // OpenAI model to be used for generating the dictionary.
	languageLevel types.LanguageLevel // Proficiency or complexity level of the language used.
	languageFrom  types.Language      // Source language for the dictionary craft.
	languageTo    types.Language      // Target language for the dictionary craft.

	temperature       float64 // Controls the randomness or creativity of the generation process.
	maxConcurrent     int     // Maximum number of concurrent workers processing the generation.
	dictionariesCount int     // Specifies how many dictionaries should be generated.
	words             int     // Desired number of words to be generated.
	dictionaries      int     // Duplicate of dictionariesCount indicating number of dictionaries to generate.
}

// NewDictionaryCraftData creates a new DictionaryCraftData instance.
func NewDictionaryCraftData() DictionaryCraftData {
	return DictionaryCraftData{
		promptBuf:         &bytes.Buffer{},
		maxConcurrent:     defaultConcurrent,
		dictionariesCount: defaultDictionaries,
	}
}

// GetDictionaryName returns the dictionary name crafted by OpenAI.
func (r *DictionaryCraftData) GetDictionaryName() string {
	if r.response != nil {
		return r.response.Meta.Name
	}
	return ""
}

// GetDictionaryOverview return the dictionary description crafted by openAI.
func (r *DictionaryCraftData) GetDictionaryOverview() string {
	if r.response != nil {
		return r.response.Meta.Description
	}
	return ""
}

// GetDictionaryAuthor return the dictionary author crafted by openAI.
func (r *DictionaryCraftData) GetDictionaryAuthor() string {
	if r.response != nil {
		return r.response.Meta.Author
	}
	return ""
}

// GetWordsContainer returns the dictionary words crafted by OpenAI.
func (r *DictionaryCraftData) GetWordsContainer() WordsContainer {
	if r.response == nil {
		return WordsContainer{Words: []DictionaryWordFromAI{}}
	}
	return WordsContainer{Words: r.response.Words}
}

// GetPrompt returns the prompt name.
func (r *DictionaryCraftData) GetPrompt() string {
	return r.prompt
}

// GetFilename returns the dictionary filename.
func (r *DictionaryCraftData) GetFilename() string {
	return r.filename
}

// GetDictionaryDescription returns the dictionary description (for AI model).
func (r *DictionaryCraftData) GetDictionaryDescription() string {
	return r.dictionaryDescription
}

// GetDictionaryTopic returns the dictionary topic (for AI model).
func (r *DictionaryCraftData) GetDictionaryTopic() string {
	return r.dictionaryTopic
}

// GetModel returns the model.
func (r *DictionaryCraftData) GetModel() chatgpt.OpenAIModel {
	return r.model
}

// GetLanguageLevel returns the language level.
func (r *DictionaryCraftData) GetLanguageLevel() types.LanguageLevel {
	return r.languageLevel
}

// GetLanguageFrom returns the source language.
func (r *DictionaryCraftData) GetLanguageFrom() types.Language {
	return r.languageFrom
}

// GetLanguageTo returns the target language.
func (r *DictionaryCraftData) GetLanguageTo() types.Language {
	return r.languageTo
}

// GetWordsCount returns the number of words to be generated.
func (r *DictionaryCraftData) GetWordsCount() int {
	return r.words
}

// GetSubcategory returns the dictionary subcategory.
func (r *DictionaryCraftData) GetSubcategory() string {
	return r.languageFrom.Code + "-" + r.languageTo.Code
}

// getPromptBody returns an io.Reader with prompt content.
func (r *DictionaryCraftData) getPromptBody() io.Reader {
	return r.promptBuf
}

// ToPromptTemplate converts the DictionaryCraftData into a promptTemplate.
func (r *DictionaryCraftData) toPromptTemplate() craftDictionaryPromptTemplate {
	return craftDictionaryPromptTemplate{
		DictionaryDescription: r.dictionaryDescription,
		DictionaryTopic:       r.dictionaryTopic,
		LanguageLevel:         r.languageLevel.String(),
		LanguageFrom:          r.languageFrom.Name,
		LanguageTo:            r.languageTo.Name,
		WordsCount:            r.words,
	}
}

// Setup initializes the DictionaryCraftData with the provided request parameters.
func (r *DictionaryCraftData) Setup(
	ctx context.Context,
	req *RequestDictionaryCraft,
	s3cli *cloud.Bucket,
	promptBucketName string,
) error {
	r.request = req.Clone()
	r.filename = uuid.New().String()

	var (
		wg       sync.WaitGroup
		results  = make(chan workerResult, 3)
		setupErr error
	)

	runWorker(ctx, &wg, results, "openai", func() error {
		return r.setupOpenAI(ctx, req, s3cli, promptBucketName)
	})
	runWorker(ctx, &wg, results, "dictionary", func() error {
		return r.setupDictionaryMeta(req)
	})
	runWorker(ctx, &wg, results, "languages", func() error {
		return r.setupLanguages(req)
	})

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.error != nil {
			if setupErr == nil {
				setupErr = errors.Join(ErrorWorkerProcess(res.key), res.error)
			} else {
				setupErr = errors.Join(setupErr, ErrorWorkerProcess(res.key), res.error)
			}
		}
	}

	if setupErr != nil {
		return errors.Join(ErrorSetupProcess, setupErr)
	}

	// fetch prompt content
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s3Response, err := s3cli.GetObjectBody(ctx, r.prompt, promptBucketName)
		if err != nil {
			return errors.Join(ErrorGetBucketFileContent(r.prompt, promptBucketName), err)
		}
		defer s3Response.Close()
		r.promptBuf.Reset()
		return utils.TemplateFromReaderToWriter(r.promptBuf, s3Response, r.toPromptTemplate())
	}
}

func (r *DictionaryCraftData) setupOpenAI(ctx context.Context, req *RequestDictionaryCraft, s3cli *cloud.Bucket, bucket string) error {
	if req.PromptName == nil {
		prompt, err := s3cli.GetRandomKey(ctx, bucket, craftPromptPrefix)
		if err != nil {
			return errors.Join(ErrorGetKeyFromBucket(craftPromptPrefix, bucket), err)
		}
		r.prompt = prompt
	} else {
		r.prompt = aws.ToString(req.PromptName)
	}

	if req.OpenaiModel == nil {
		r.model = defaultModel
	} else {
		model, err := chatgpt.ParseModel(aws.ToString(req.OpenaiModel))
		if err != nil {
			return errors.Join(ErrorOpenAIModelNotSupported(aws.ToString(req.OpenaiModel)), err)
		}
		r.model = model
	}

	if req.Temperature != nil {
		val := aws.ToFloat64(req.Temperature)
		if val < 0 || val > 1 {
			return ErrorGetTemperature(val, 0, 1)
		}
		r.temperature = val
	} else {
		r.temperature = defaultTemperature
	}

	if req.MaxConcurrent != nil && aws.ToInt(req.MaxConcurrent) >= 1 {
		r.maxConcurrent = aws.ToInt(req.MaxConcurrent)
	} else {
		r.maxConcurrent = defaultConcurrent
	}

	return nil
}

func (r *DictionaryCraftData) setupDictionaryMeta(req *RequestDictionaryCraft) error {
	if req.DictionaryTopic != nil {
		r.dictionaryTopic = aws.ToString(req.DictionaryTopic)
	} else {
		topic, err := types.GetRandomDictionaryTopic()
		if err != nil {
			return errors.Join(ErrorGenerateDictionaryTopic, err)
		}
		r.dictionaryTopic = topic.String()
	}

	if req.DictionaryDescription != nil {
		r.dictionaryDescription = aws.ToString(req.DictionaryDescription)
	} else {
		description, err := types.GetRandomDictionaryDescription()
		if err != nil {
			return errors.Join(ErrorGenerateDictionaryDescription, err)
		}
		r.dictionaryDescription = description.String()
	}

	if req.WordsCount != nil {
		val := aws.ToInt(req.WordsCount)
		if val < 1 || val > dictionaryMaxLength {
			return ErrorGetWordsCount(val, 1, dictionaryMaxLength)
		}
		r.words = val
	} else {
		val, err := utils.RandomInt(dictionaryMinLength, dictionaryMaxLength)
		if err != nil {
			return errors.Join(ErrorGenerateRandomInt(dictionaryMinLength, dictionaryMaxLength), err)
		}
		r.words = val
	}

	if req.DictionariesCount != nil && aws.ToInt(req.DictionariesCount) >= 1 {
		r.dictionaries = aws.ToInt(req.DictionariesCount)
	} else {
		r.dictionaries = 1
	}

	return nil
}

func (r *DictionaryCraftData) setupLanguages(req *RequestDictionaryCraft) error {
	var err error

	if req.LanguageLevel != nil {
		r.languageLevel, err = types.ParseLanguageLevel(aws.ToString(req.LanguageLevel))
	} else {
		r.languageLevel, err = types.GetRandomLanguageLevel()
	}
	if err != nil {
		return errors.Join(ErrorGenerateLanguage("level"), err)
	}

	if req.LanguageFrom != nil {
		r.languageFrom, err = types.ParseLanguageString(aws.ToString(req.LanguageFrom))
	} else {
		r.languageFrom, err = types.GetRandomLanguage()
	}
	if err != nil {
		return errors.Join(ErrorGenerateLanguage("from"), err)
	}

	if req.LanguageTo != nil {
		langTo, err := types.ParseLanguageString(aws.ToString(req.LanguageTo))
		if err == nil && langTo.Code != r.languageFrom.Code {
			r.languageTo = langTo
			return nil
		}
	}
	r.languageTo, err = types.GetRandomLanguageExcept(r.languageFrom)
	if err != nil {
		return errors.Join(ErrorGenerateLanguage("to"), err)
	}
	return nil
}
