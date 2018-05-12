package kafka

import (
	"github.com/arizz96/event-api/config"

	"encoding/json"

	"github.com/arizz96/event-api/logging"
	v1types "github.com/arizz96/event-api/types/v1"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Writer puts messages into a kafka topic
type Writer struct {
	producer *kafka.Producer
	topic    string
}

const (
	flushTimeout = 10000
)

var (
	// Prometheus metrics
	txBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_tx_bytes_total",
		Help: "The number of bytes being produced to kafka brokers",
	}, []string{"broker", "producer"})

	rxBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_rx_bytes_total",
		Help: "The number of bytes being received to kafka brokers",
	}, []string{"broker", "producer"})

	txRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_tx_requests_total",
		Help: "The number of requests sent to brokers",
	}, []string{"broker", "producer"})

	rxResponses = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_rx_responses_total",
		Help: "The number of responses from the brokers",
	}, []string{"broker", "producer"})

	latency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "kafka_producer_latency_ms",
		Help: "Average request latency",
	}, []string{"broker", "producer"})
)

func init() {
	prometheus.MustRegister(txBytes)
	prometheus.MustRegister(rxBytes)
	prometheus.MustRegister(txRequests)
	prometheus.MustRegister(rxResponses)
	prometheus.MustRegister(latency)
}

// NewWriter creates and returns a new Kafka Writer
func NewWriter() *Writer {
	log := logging.GetLogger(logrus.Fields{"package": "kafka"})

	// Set up Kafka producer config
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        config.AppConfig.KafkaBrokers,
		"statistics.interval.ms":   500,
		"request.required.acks":    -1,
		"message.timeout.ms":       50000,
		"queue.buffering.max.ms":   5000,
		"message.send.max.retries": 10,
	}

	// Set Kafka debugging if in DebugMode
	if config.AppConfig.DebugMode == true {
		cfg.SetKey("debug", "protocol,topic,msg")
	}

	// Create a new kafka producer
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create kafka client: %s\n", err)
	}

	// Create and return a new Kafka Writer
	writer := &Writer{
		producer: producer,
		topic:    config.AppConfig.KafkaTopic,
	}

	// Fork off to handle Kafka events
	go writer.handleEvents()

	// Return the Kafka Writer
	return writer
}

// Write writes a given event to a configured topic in the kafka cluster
func (h *Writer) Write(ev v1types.Message) error {
	// log := logging.GetLogger(logrus.Fields{"package": "kafka"})

	// Create message
	h.producer.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &h.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(ev.GetMessageID()),
		Value: []byte(ev.String()),
	}

	return nil
}

func (h *Writer) handleEvents() {
	log := logging.GetLogger(logrus.Fields{"package": "kafka"})

	// Loop over events
	for e := range h.producer.Events() {
		switch ev := e.(type) {

		// Delivery reports
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Errorf("Delivery failed: %v", ev.TopicPartition.Error)
			} else {
				log.Debugf("Delivered message to %v\n", ev.TopicPartition)
			}

		// Stats updates
		case *kafka.Stats:
			// Unmarshal Stats object
			var stats Stats
			err := json.Unmarshal([]byte(e.String()), &stats)
			if err != nil {
				log.Errorf("Error unmarshalling Kafka Stats: %s", err)
			}

			// Loop over brokers
			for _, v := range stats.Brokers {
				// Create common labels
				lbls := prometheus.Labels{"broker": v.Name, "producer": "event-api"}

				// Record metrics
				txBytes.With(lbls).Set(float64(v.Txbytes)) // Total bytes sent to broker
				rxBytes.With(lbls).Set(float64(v.Rxbytes)) // Total bytes received by broker
				txRequests.With(lbls).Set(float64(v.Tx))   // Total requests sent
				rxResponses.With(lbls).Set(float64(v.Rx))  // Total responses received
				latency.With(lbls).Set(float64(v.Rtt.Avg)) // Avg roundtrip to broker
			}
		}
	}
}

// Close cleans up the Kafka Writer
func (h *Writer) Close() {
	log := logging.GetLogger(logrus.Fields{"package": "kafka"})
	log.Info("Shutting down producer")

	// Flush any remaining messages, waiting up until the flushTimeout
	msgs := h.producer.Flush(flushTimeout)
	if msgs != 0 {
		log.Errorf("Failed to flush %d messages after %d ms", msgs, flushTimeout)
	} else {
		log.Info("All messages have been flushed")
	}

	// Defer close until function exists
	h.producer.Close()
	log.Info("Producer has been closed")
}
