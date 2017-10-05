package kinesis

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/sirupsen/logrus"
)

type IngestionHandler struct {
	kc         kinesisiface.KinesisAPI
	streamName *string
	log        *logrus.Logger
}

func NewIngestionHandler(log *logrus.Logger) *IngestionHandler {
	logger := log.WithFields(logrus.Fields{"package": "kinesis", "function": "NewIngestionHandler"})
	s, err := session.NewSession()
	if err != nil {
		logger.Fatal(err)
	}
	h := &IngestionHandler{
		kc:  kinesis.New(s),
		log: log,
	}
	h.streamName = aws.String(config.Get().StreamName)
	return h
}

func NewMockLocalStackIngestionHandler(log *logrus.Logger) *IngestionHandler {
	logger := log.WithFields(logrus.Fields{"package": "kinesis", "function": "NewMockLocalStackIngestionHandler"})
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		logger.Fatal(err)
	}
	h := &IngestionHandler{
		kc: kinesis.New(s, &aws.Config{
			Endpoint: aws.String("http://192.168.1.225:4568"),
		}),
		log: log,
	}
	h.streamName = aws.String(config.Get().StreamName)
	return h
}

func NewMockIngestionHandler() *IngestionHandler {
	return &IngestionHandler{
		kc: NewMockKinesisClient(),
	}
}

func (h *IngestionHandler) Start() error {
	return nil
}

func (h *IngestionHandler) Shutdown() error {
	return nil
}

func (h *IngestionHandler) ProcessMessage(r, partition string) {
	logger := h.log.WithFields(logrus.Fields{"package": "kinesis", "function": "ProcessMessage"})
	_, err := h.kc.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte(r),
		PartitionKey: aws.String(partition),
		StreamName:   h.streamName,
	})
	if err != nil {
		logger.Info(err)
		return
	}
}
