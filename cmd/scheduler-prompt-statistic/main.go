// Package main implements a Lambda function go handle dictionary craft statistics.
package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/tg"

	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	defaultMaxWorkers = 1
	telegramChatID    = -100
)

var (
	awsRegion      = os.Getenv("AWS_REGION")
	telegramToken  = os.Getenv("TELEGRAM_TOKEN")
	dbDynamo       *cloud.Dynamo
	telegramClient tg.Telegram
)

func init() {
	debug.SetGCPercent(500)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	dbDynamo = cloud.NewDynamo(cfg)

	if telegramToken != "" {
		telegramClient, err = tg.New(telegramToken, telegramChatID)
		if err != nil {
			log.Printf("Error initializing Telegram client: %v", err)
		}
	}
}

// func handler(ctx context.Context, _ zerolog.Logger, _ json.RawMessage) error {
// 	items, err := fetchProcessingItems(ctx, dbDynamo, applingoprocessing.TableName)
// 	if err != nil {
// 		return err
// 	}

// 	posts := GeneratePosts(generateModelStats(items))

// 	return nil
// }

func main() {

	items, err := fetchProcessingItems(context.TODO(), dbDynamo, applingoprocessing.TableName)
	if err != nil {
		panic(err)
	}

	posts := GeneratePosts(generateModelStats(items))
	for _, post := range posts {
		_, err := telegramClient.Send(post.Content, tg.ModeMarkdown)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)

	}

	// lambda.Start(
	// 	trigger.NewLambda(
	// 		trigger.Config{
	// 			MaxWorkers: defaultMaxWorkers,
	// 		},
	// 		handler,
	// 	).Handle,
	// )
}
