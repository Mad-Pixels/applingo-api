package forge

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// checkDictionaryPromptTemplate defines the structure for the prompt used in dictionary check.
// It aggregates the main parameters that will be injected into the AI prompt template.
type checkDictionaryPromptTemplate struct {
	// DictionaryDescription is a brief overview of the dictionary's content for the AI model.
	DictionaryDescription string
	// DictionaryTopic specifies the main subject or theme of the dictionary.
	DictionaryTopic string
	// DictionaryName of the dictionary
	DictionaryName string
	// DictionaryOverview a visible dictionary description.
	DictionaryOverview string
	// LanguageLevel represents the CEFR proficiency level as a string (e.g., "A1", "B2").
	LanguageLevel string
	// LanguageFrom indicates the source language in string format.
	LanguageFrom string
	// LanguageTo indicates the target language in string format.
	LanguageTo string
}

// DictionaryCheckData holds all data related to a dictionary check request.
// It encapsulates both the input parameters and internal data used during the check process.
type DictionaryCheckData struct {
	request       *RequestDictionaryCheck        // Original dictionary check request parameters.
	response      *ResponseDictionaryCheck       // Original dictionary check response data.
	item          *applingoprocessing.SchemaItem // DynamoDB processing table row with data.
	dictionaryBuf *bytes.Buffer                  // Internal buffer to store the fetched dictionary content.
	promptBuf     *bytes.Buffer                  // Internal buffer to store the fetched prompt template.

	prompt      string              // Prompt string used for initiating dictionary generation.
	model       chatgpt.OpenAIModel // OpenAI model to be used for generating the dictionary.
	temperature float64             // Controls the randomness or creativity of the generation process.
}

// NewDictionaryCheckData creates a new DictionaryCheckData instance.
func NewDictionaryCheckData() DictionaryCheckData {
	return DictionaryCheckData{
		dictionaryBuf: &bytes.Buffer{},
		promptBuf:     &bytes.Buffer{},
	}
}

// GetScore return the dictionary score num.
func (r *DictionaryCheckData) GetScore() int {
	if r.response != nil {
		return r.response.Meta.Score
	}
	return 0
}

// GetReason return a score description message.
func (r *DictionaryCheckData) GetReason() string {
	if r.response != nil {
		return r.response.Meta.Reason
	}
	return ""
}

// GetPrompt returns the prompt name.
func (r *DictionaryCheckData) GetPrompt() string {
	return r.prompt
}

// GetModel returns the model.
func (r *DictionaryCheckData) GetModel() chatgpt.OpenAIModel {
	return r.model
}

// GetTemperature returns the temperature.
func (r *DictionaryCheckData) GetTemperature() float64 {
	return r.temperature
}

// getPromptBody returns an io.Reader with prompt content.
func (r *DictionaryCheckData) getPromptBody() io.Reader {
	return io.MultiReader(r.promptBuf, bytes.NewBufferString("\n\n"), r.dictionaryBuf)
}

// ToPromptTemplate converts the DictionaryCheckData into a promptTemplate.
func (r *DictionaryCheckData) toPromptTemplate() checkDictionaryPromptTemplate {
	return checkDictionaryPromptTemplate{
		DictionaryDescription: r.item.Description,
		DictionaryTopic:       r.item.Topic,
		DictionaryName:        r.item.Name,
		DictionaryOverview:    r.item.Overview,
		LanguageLevel:         r.item.Level,
		LanguageFrom:          utils.SplitValues(r.item.Languages)[0],
		LanguageTo:            utils.SplitValues(r.item.Languages)[1],
	}
}

// Setup initializes the DictionaryCheckData with the provided request parameters.
func (r *DictionaryCheckData) Setup(
	ctx context.Context,
	req *RequestDictionaryCheck,
	item *applingoprocessing.SchemaItem,
	s3cli *cloud.Bucket,
	promptBucketName string,
	dictionaryBucketName string,
) error {
	var (
		results  = make(chan workerResult, 2)
		setupErr error
		wg       sync.WaitGroup
	)
	r.request = req.Clone()
	r.item = item

	// prompt.
	if req.PromptName == nil {
		prompt, err := s3cli.GetRandomKey(ctx, promptBucketName, checkPromptPrefix)
		if err != nil {
			return errors.Join(ErrorGetKeyFromBucket(checkPromptPrefix, promptBucketName), err)
		}
		r.prompt = prompt
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

	// Dictionary worker.
	runWorker(ctx, &wg, results, "dictionary", func() error {
		r.dictionaryBuf.Reset()
		if err := s3cli.Read(ctx, r.dictionaryBuf, utils.RecordToFileID(r.item.Id), dictionaryBucketName); err != nil {
			return errors.Join(ErrorGetBucketFileContent(r.item.Id, dictionaryBucketName), err)
		}
		return nil
	})

	// Prompt worker.
	runWorker(ctx, &wg, results, "prompt", func() error {
		resp, err := s3cli.GetObjectBody(ctx, r.prompt, promptBucketName)
		if err != nil {
			return errors.Join(ErrorGetBucketFileContent(r.prompt, promptBucketName), err)
		}
		defer resp.Close()

		r.promptBuf.Reset()
		if err = utils.TemplateFromReaderToWriter(r.promptBuf, resp, r.toPromptTemplate()); err != nil {
			return errors.Join(ErrorParseTemplate(r.prompt), err)
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
	return nil
}
