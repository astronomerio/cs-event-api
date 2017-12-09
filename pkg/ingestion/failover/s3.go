package failover

import (
	"bytes"
	"context"
	"github.com/astronomerio/event-api/pkg/config"
	"github.com/astronomerio/event-api/pkg/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"time"
)

func init() {
	c = config.Get()
}

var (
	c       *config.Configuration
	timeout time.Duration
)

func UploadMessage(msg kafka.Message) {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "ProcessMessage"})
	logger.Debug("uploading message to s3 due to kafka outage")

	sess := session.Must(session.NewSession(
		&aws.Config{
			Region:      &c.S3Region,
			Credentials: credentials.NewEnvCredentials(),
		}))
	svc := s3.New(sess)

	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
		// Ensure the context is canceled to prevent leaking.
		// See context package for more information, https://golang.org/pkg/context/
		defer cancelFn()
	}
	msgLen := int64(len(msg.Value))
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.S3Bucket),
		Key:           aws.String(string(msg.Key)),
		Body:          bytes.NewReader(msg.Value),
		ContentLength: &msgLen,
	})
	if err != nil {
		logger.Errorf("error uploading to s3 %s", err.Error())
	}
}
