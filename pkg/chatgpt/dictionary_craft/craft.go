package dictionary_craft

import (
	"context"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"

	"github.com/pkg/errors"
)

// Craft generates a dictionary using the ChatGPT API.
func Craft(ctx context.Context, req *Request, promptBucket string, chatgptCli *chatgpt.Client, s3Cli *cloud.Bucket) (*Dictionary, error) {
	if err := req.Prepare(ctx, s3Cli, promptBucket); err != nil {
		return nil, errors.Wrap(err, "failed to prepare request data")
	}
	gptReq := chatgpt.NewRequest(
		req.GetModel(),
		[]chatgpt.Message{chatgpt.NewUserMessage(req.GetPromptBody())},
	)
	resp, err := chatgptCli.SendMessage(ctx, gptReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process ChatGPT request")
	}

	var dictionary Dictionary
	if err := dictionary.Unmarshal([]byte(resp.GetResponseText())); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal dictionary")
	}
	if len(dictionary.Words) == 0 {
		return nil, errors.New("dictionary has no words")
	}

	dictionary.Meta = *req
	dictionary.Meta.DictionaryLength = len(dictionary.Words)
	return &dictionary, nil
}
