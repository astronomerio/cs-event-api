package kafka

import (
	"strings"

	"github.com/astronomerio/event-api/config"

	"encoding/json"
	"errors"
	"fmt"

	"github.com/astronomerio/event-api/logging"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Writer puts messages into a kafka topic
type Writer struct {
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

func NewWriter() *Writer {
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "NewWriter"})
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(appConfig.KafkaBrokers, ","),
		"statistics.interval.ms":   500,
		"request.required.acks":    -1,
		"message.timeout.ms":       50000,
		"queue.buffering.max.ms":   5000,
		"message.send.max.retries": 10,
	}
	if config.Get().DebugMode == true {
		cfg.SetKey("debug", "protocol,topic,msg")
	}

	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("failed to create kafka client: %s\n", err)
	}
	h := Writer{
		producer: producer,
		topic:    appConfig.KafkaTopic,
	}
	return &h
}

func (h *Writer) Start() error {
	h.startEventListener()
	return nil
}

func (h *Writer) Shutdown() error {
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "Shutdown"})
	log.Info("shutting down Kafka handler")
	defer h.producer.Close()

	msgs := h.producer.Flush(10000)
	if len(h.producer.ProduceChannel()) != 0 {
		return errors.New(fmt.Sprintf("%d messages were not flushed after a timeout of %d", msgs, 10000))
	}

	return nil
}

func (h *Writer) startEventListener() {
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "startEventListener"})

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
					log.Errorf("json unmarshal error: %s", err)
				}
				for _, v := range stats.Brokers {
					bytesOut.With(prometheus.Labels{"broker": v.Name, "producer": "event-api"}).Set(float64(v.Rxbytes))
					latency.With(prometheus.Labels{"broker": v.Name, "producer": "event-api"}).Set(float64(v.Rtt.Avg))
					responseRate.With(prometheus.Labels{"broker": v.Name, "producer": "event-api"}).Set(float64(v.Rx))
					requestRate.With(prometheus.Labels{"broker": v.Name, "producer": "event-api"}).Set(float64(v.Tx))
				}
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Errorf("delivery failed: %v", m.TopicPartition.Error)
				}
			default:
				log.Errorf("non kafka message found in event stream: %s\n", ev)
			}
		}
	}()
}

func (h *Writer) ProcessMessage(message, partition string) {
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka", "function": "ProcessMessage"})

	// Create message
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &h.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(partition),
		Value: []byte(message),
	}

	// Skip if we don't have a producer running
	if isRunning != true {
		log.Error("Event listener isn't active")
		return
	}

	// Produce the events
	err := h.producer.Produce(msg, h.producer.Events())
	if err != nil {
		log.Errorf("Unable to produce %f", err.Error())
	}

}
