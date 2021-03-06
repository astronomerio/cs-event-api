package kafka

import (
	"strings"

	"github.com/astronomerio/event-api/pkg/config"

	"encoding/json"
	"errors"
	"fmt"
	"github.com/astronomerio/event-api/pkg/ingestion/failover"
	"github.com/astronomerio/event-api/pkg/logging"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type KafkaHandler struct {
	producer *kafka.Producer
	topic    string
}

var (
	bytesOut = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_message_out_bytes_total",
		Help: "The number of bytes being produced to kafka brokers",
	}, []string{"broker", "producer"})
	requestRate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_requests_total",
		Help: "Average number of requests",
	}, []string{"broker", "producer"})
	responseRate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_responses_total",
		Help: "Average number of responses received",
	}, []string{"broker", "producer"})
	latency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_latency_ms",
		Help: "Average request latency",
	}, []string{"broker", "producer"})
)

func init() {
	prometheus.MustRegister(bytesOut)
	prometheus.MustRegister(requestRate)
	prometheus.MustRegister(responseRate)
	prometheus.MustRegister(latency)
}

var appConfig = config.Get()
var isRunning = false

func NewHandler() *KafkaHandler {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "NewHandler"})
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":      strings.Join(appConfig.KafkaBrokers, ","),
		"statistics.interval.ms": 500,
		"request.required.acks":  -1,
		"message.timeout.ms":     50000,
		"queue.buffering.max.ms": 5000,
		"message.send.max.retries": 10,
	}
	if config.Get().DebugMode == true {
		cfg.SetKey("debug", "protocol,topic,msg")
	}

	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		logger.Fatalf("failed to create kafka client: %s\n", err)
	}
	h := KafkaHandler{
		producer: producer,
		topic:    appConfig.KafkaTopic,
	}
	return &h
}

func (h *KafkaHandler) Start() error {
	h.startEventListener()
	return nil
}

func (h *KafkaHandler) Shutdown() error {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "Shutdown"})
	logger.Info("shutting down Kafka handler")
	defer h.producer.Close()

	msgs := h.producer.Flush(10000)
	if len(h.producer.ProduceChannel()) != 0 {
		return errors.New(fmt.Sprintf("%d messages were not flushed after a timeout of %d", msgs, 10000))
	}

	return nil
}

func (h *KafkaHandler) startEventListener() {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "startEventListener"})
	go func() {
		isRunning = true
		defer func() {
			isRunning = false
		}()
		for e := range h.producer.Events() {
			switch ev := e.(type) {
			case *kafka.Stats:
				var stats Stats
				err := json.Unmarshal([]byte(e.String()), &stats)
				if err != nil {
					logger.Errorf("json unmarshal error: %s", err)
				}
				for _, v := range stats.Brokers {
					bytesOut.With(prometheus.Labels{"broker": v.Name, "producer": "ingestion-api"}).Set(float64(v.Rxbytes))
					latency.With(prometheus.Labels{"broker": v.Name, "producer": "ingestion-api"}).Set(float64(v.Rtt.Avg))
					responseRate.With(prometheus.Labels{"broker": v.Name, "producer": "ingestion-api"}).Set(float64(v.Rx))
					requestRate.With(prometheus.Labels{"broker": v.Name, "producer": "ingestion-api"}).Set(float64(v.Tx))
				}
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					logger.Errorf("delivery failed: %v", m.TopicPartition.Error)
					if config.Get().EnableFailover == true {
						failover.UploadMessage(*m)
					}
				}
			default:
				logger.Errorf("non kafka message found in event stream: %s\n", ev)
			}
		}
	}()
}

func (h *KafkaHandler) ProcessMessage(message, partition string) {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "ProcessMessage"})
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &h.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(partition),
		Value: []byte(message),
	}

	if isRunning != true {
		logger.Error("event listener isn't active")
	} else {
		err := h.producer.Produce(msg, h.producer.Events())
		if err != nil {
			logger.Errorf("unable to produce %f", err.Error())
		}
	}
}
