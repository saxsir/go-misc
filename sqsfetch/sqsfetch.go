package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var version = "ビルド時に埋め込まれます"

func main() {
	var (
		endpoint    = flag.String("endpoint", "", "optional aws endpoint url")
		queueURL    = flag.String("queue-url", "", "optional aws endpoint url")
		maxNumber   = flag.String("n", 1, "")
		region      = flag.String("region", "ap-northeast-1", "aws regions")
		showVersion = flag.Bool("version", false, "prints version.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stdout, "version: %s\n", version)
		return
	}

	// sqsから読み込む準備
	sess, err := session.NewSession(&aws.Config{
		// 指定しなければaws-sdkのデフォルトになります
		Endpoint: aws.String(*endpoint),
		Region:   aws.String(*region),
	})
	if err != nil {
		log.Fatal(err)
	}

	svc := sqs.New(sess)
	res, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MaxNumberOfMessages: aws.Int64(maxNumber),
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:          queueURL,
		VisibilityTimeout: aws.Int64(0),
		WaitTimeSeconds:   aws.Int64(0),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range res.Messages {
		fmt.Printf("%#v\n", m)
	}
}
