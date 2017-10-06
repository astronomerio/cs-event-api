package kinesis

import (
	"github.com/astronomerio/clickstream-ingestion-api/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/sirupsen/logrus"
	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
)

type Handler struct {
	kc         kinesisiface.KinesisAPI
	streamName *string
}

func NewHandler() *Handler {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kinesis", "function": "NewHandler"})
	s, err := session.NewSession()
	if err != nil {
		logger.Fatal(err)
	}
	h := &Handler{
		kc:  kinesis.New(s),
	}
	h.streamName = aws.String(config.Get().StreamName)
	return h
}

func NewMockLocalStackHandler() *Handler {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kinesis", "function": "NewMockLocalStackHandler"})
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		logger.Fatal(err)
	}
	h := &Handler{
		kc: kinesis.New(s, &aws.Config{
			Endpoint: aws.String("http://192.168.1.225:4568"),
		}),
	}
	h.streamName = aws.String(config.Get().StreamName)
	return h
}

func NewMockHandler() *Handler {
	return &Handler{
		kc: NewMockKinesisClient(),
	}
}

func (h *Handler) Start() error {
	return nil
}

func (h *Handler) Shutdown() error {
	return nil
}

func (h *Handler) ProcessMessage(r, partition string) {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kinesis", "function": "ProcessMessage"})
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
