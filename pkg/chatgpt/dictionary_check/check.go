package dictionary_check

import (
	"context"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/pkg/errors"
)

// Check checks the dictionary.
func Check(ctx context.Context, req *Request, promptBucket, dictionaryBucket string, chatgptCli *chatgpt.Client, s3Cli *cloud.Bucket) (*Result, error) {
	if err := req.Prepare(ctx, s3Cli, promptBucket, dictionaryBucket); err != nil {
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

	var result Result
	if err := result.Unmarshal([]byte(resp.GetResponseText())); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal result")
	}
	return &result, nil
}
