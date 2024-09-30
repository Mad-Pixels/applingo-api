package main

import (
	"context"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/Mad-Pixels/lingocards-api/pkg/logger"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"runtime/debug"
	"sync"
)

var (
	servicePutScvQueueUrl = os.Getenv("SERVICE_PUT_CSV_QUEUE_URL")
	awsRegion             = os.Getenv("AWS_REGION")

	sqsQueue *cloud.Queue
	log      zerolog.Logger
)

func init() {
	debug.SetGCPercent(500)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	sqsQueue = cloud.NewQueue(cfg)
	log = logger.InitLogger()
}

func handler(ctx context.Context, e events.DynamoDBEvent) error {
	var (
		wg   sync.WaitGroup
		errs = make(chan error, len(e.Records))
	)
	for _, eventRecord := range e.Records {
		wg.Add(1)

		go func(e events.DynamoDBEventRecord) {
			defer wg.Done()

			payload, err := serializer.MarshalJSON(e)
			if err != nil {
				errs <- err
				return
			}
			_, err = sqsQueue.SendMessage(ctx, servicePutScvQueueUrl, string(payload))
			if err != nil {
				errs <- err
				return
			}
		}(eventRecord)
	}
	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			log.Error().Err(err).Msg("Error processing record")
		}
	}
	if len(errs) > 0 {
		return errors.New("finish with errors")
	}
	return nil
}

func main() { lambda.Start(handler) }
