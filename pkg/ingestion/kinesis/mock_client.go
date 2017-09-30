package kinesis

import (
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type mockKinesisClient struct {
	kinesisiface.KinesisAPI
}

func NewMockKinesisClient() *mockKinesisClient {
	return &mockKinesisClient{}
}

func (h *mockKinesisClient) Start() error {
	return nil
}

func (h *mockKinesisClient) Shutdown() error {
	return nil
}

func (m *mockKinesisClient) PutRecord(*kinesis.PutRecordInput) (*kinesis.PutRecordOutput, error) {
	return nil, nil
}
