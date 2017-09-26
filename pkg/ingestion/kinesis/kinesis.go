package kinesis

import (
	"log"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type KinesisIngestionHandler struct {
	kc         kinesisiface.KinesisAPI
	streamName *string
}

func NewKinesisIngestionHandler() *KinesisIngestionHandler {
	s, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	h := &KinesisIngestionHandler{
		kc: kinesis.New(s),
	}
	h.streamName = aws.String(config.Get().StreamName)
	return h
}

func NewMockKinesisLocalStackIngestionHandler() *KinesisIngestionHandler {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		log.Fatal(err)
	}
	h := &KinesisIngestionHandler{
		kc: kinesis.New(s, &aws.Config{
			Endpoint: aws.String("http://192.168.1.225:4568"),
		}),
	}
	h.streamName = aws.String(config.Get().StreamName)
	return h
}

func NewMockKinesisIngestionHandler() *KinesisIngestionHandler {
	return &KinesisIngestionHandler{
		kc: NewMockKinesisClient(),
	}
}

func (h *KinesisIngestionHandler) ProcessMessage(r, partition string) {
	_, err := h.kc.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte(r),
		PartitionKey: aws.String(partition),
		StreamName:   h.streamName,
	})
	if err != nil {
		log.Println(err)
		return
	}
}
