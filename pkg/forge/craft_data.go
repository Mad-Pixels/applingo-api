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

// GetDicionaryName return the name crafted by openAI.
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

// GetWords return the dictionary words crafted by openAI.
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
func (d *DictionaryCraftData) toPromptTemplate() craftDictionaryPromptTemplate {
	return craftDictionaryPromptTemplate{
		DictionaryDescription: d.dictionaryDescription,
		DictionaryTopic:       d.dictionaryTopic,
		LanguageLevel:         d.languageLevel.String(),
		LanguageFrom:          d.languageFrom.Name,
		LanguageTo:            d.languageTo.Name,
		WordsCount:            d.words,
	}
}

// Setup initializes the DictionaryCraftData with the provided request parameters.
func (r *DictionaryCraftData) Setup(ctx context.Context, req *RequestDictionaryCraft, s3cli *cloud.Bucket, promptBucketName string) error {
	var (
		results  = make(chan workerResult, 3)
		setupErr error
		wg       sync.WaitGroup
	)
	r.request = req.Clone()

	// OpenAI model data defenition worker.
	runWorker(ctx, &wg, results, "openai", func() (err error) {
		// prompt.
		if req.PromptName == nil {
			prompt, err := s3cli.GetRandomKey(ctx, promptBucketName, craftPromptPrefix)
			if err != nil {
				return errors.Join(ErrorGetKeyFromBucket(craftPromptPrefix, promptBucketName), err)
			}
			r.prompt = prompt
		} else {
			r.prompt = aws.ToString(req.PromptName)
		}
		// model.
		if req.OpenaiModel == nil {
			r.model = defaultModel
		} else {
			model, err := chatgpt.ParseModel(aws.ToString(req.OpenaiModel))
			if err != nil {
				return errors.Join(ErrorOpenAIModelNotSupported(aws.ToString(req.OpenaiModel)), err)
			}
			r.model = model
		}
		// temperature.
		if req.Temperature != nil {
			if aws.ToFloat64(req.Temperature) < 0 || aws.ToFloat64(req.Temperature) > 1 {
				return ErrorGetTemperature(aws.ToFloat64(req.Temperature), 0, 1)
			}
			r.temperature = aws.ToFloat64(req.Temperature)
		} else {
			r.temperature = defaultTemperature
		}
		// Max concurrent execution.
		if req.MaxConcurrent != nil {
			if aws.ToInt(req.MaxConcurrent) < 1 {
				r.maxConcurrent = defaultConcurrent
			}
			r.maxConcurrent = aws.ToInt(req.MaxConcurrent)
		} else {
			r.maxConcurrent = defaultConcurrent
		}
		return nil
	})

	// Dictionary data defenition worker.
	runWorker(ctx, &wg, results, "dictionary", func() (err error) {
		// Filename.
		r.filename = uuid.New().String()
		// Dictionary topic.
		if req.DictionaryTopic != nil {
			r.dictionaryTopic = aws.ToString(req.DictionaryTopic)
		} else {
			topic, err := types.GetRandomDictionaryTopic()
			if err != nil {
				return errors.Join(ErrorGenerateDictionaryTopic, err)
			}
			r.dictionaryTopic = topic.String()
		}
		// Dictionary description.
		if req.DictionaryDescription != nil {
			r.dictionaryDescription = aws.ToString(req.DictionaryDescription)
		} else {
			description, err := types.GetRandomDictionaryDescription()
			if err != nil {
				return errors.Join(ErrorGenerateDictionaryDescription, err)
			}
			r.dictionaryDescription = description.String()
		}
		// Words count.
		if req.WordsCount != nil {
			if aws.ToInt(req.WordsCount) < 1 || aws.ToInt(req.WordsCount) > dictionaryMaxLength {
				return ErrorGetWordsCount(aws.ToInt(req.WordsCount), 1, dictionaryMaxLength)
			}
			r.words = aws.ToInt(req.WordsCount)
		} else {
			words, err := utils.RandomInt(dictionaryMinLength, dictionaryMaxLength)
			if err != nil {
				return errors.Join(ErrorGenerateRandomInt(dictionaryMinLength, dictionaryMaxLength), err)
			}
			r.words = words
		}
		// Dictionaries count.
		if req.DictionariesCount != nil {
			if aws.ToInt(req.DictionariesCount) < 1 {
				r.dictionaries = 1
			}
			r.dictionaries = aws.ToInt(req.DictionariesCount)
		} else {
			r.dictionaries = 1
		}
		return nil
	})

	// Languages data defenition worker.
	runWorker(ctx, &wg, results, "languages", func() (err error) {
		// Language level.
		r.languageLevel, err = func() (types.LanguageLevel, error) {
			if req.LanguageLevel != nil {
				if level, err := types.ParseLanguageLevel(aws.ToString(req.LanguageLevel)); err == nil {
					return level, nil
				}
			}
			return types.GetRandomLanguageLevel()
		}()
		if err != nil {
			return errors.Join(ErrorGenerateLanguage("level"), err)
		}
		// Language from.
		r.languageFrom, err = func() (types.Language, error) {
			if req.LanguageFrom != nil {
				if language, err := types.ParseLanguageString(aws.ToString(req.LanguageFrom)); err == nil {
					return language, nil
				}
			}
			return types.GetRandomLanguage()
		}()
		if err != nil {
			return errors.Join(ErrorGenerateLanguage("from"), err)
		}
		// Language to.
		r.languageTo, err = func() (types.Language, error) {
			if req.LanguageTo != nil {
				if language, err := types.ParseLanguageString(aws.ToString(req.LanguageTo)); err == nil && language.Code != r.languageFrom.Code {
					return language, nil
				}
			}
			return types.GetRandomLanguageExcept(r.languageFrom)
		}()
		if err != nil {
			return errors.Join(ErrorGenerateLanguage("to"), err)
		}
		return nil
	})

	go func() {
		wg.Wait()
		close(results)
	}()
	for res := range results {
		if res.error != nil {
			workerErr := errors.Join(ErrorWorkerProcess(res.key), res.error)
			if setupErr == nil {
				setupErr = workerErr
			} else {
				setupErr = errors.Join(setupErr, workerErr)
			}
		}
	}
	if setupErr != nil {
		return errors.Join(ErrorSetupProcess, setupErr)
	}

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
		if err = utils.TemplateFromReaderToWriter(r.promptBuf, s3Response, r.toPromptTemplate()); err != nil {
			return errors.Join(ErrorParseTemplate(r.prompt), err)
		}
		return nil
	}
}
