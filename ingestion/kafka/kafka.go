package kafka

import (
	"strings"

	"github.com/astronomerio/event-api/config"

	"encoding/json"
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

const (
	flushTimeout = 10000
)

var (
	// Writer configuration
	appConfig = config.Get()
	isRunning = false

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
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka"})

	// Set up Kafka producer config
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(appConfig.KafkaBrokers, ","),
		"statistics.interval.ms":   500,
		"request.required.acks":    -1,
		"message.timeout.ms":       50000,
		"queue.buffering.max.ms":   5000,
		"message.send.max.retries": 10,
	}

	// Set Kafka debugging if in DebugMode
	if config.Get().DebugMode == true {
		cfg.SetKey("debug", "protocol,topic,msg")
	}

	// Create a new kafka producer
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create kafka client: %s\n", err)
	}

	// Create and return a new Kafka Writer
	return &Writer{
		producer: producer,
		topic:    appConfig.KafkaTopic,
	}
}

// Start initializes the kafka event listener
func (h *Writer) Start() error {
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka"})

	go func() {
		// Set our running flag
		isRunning = true

		// When this function ends, set our flag back to false
		defer func() {
			isRunning = false
		}()

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

			// Unknown message type
			default:
				log.Debugf("Received unknown event in %s\n", ev)
			}
		}
	}()

	return nil
}

// ProcessMessage writes a given message to a given topic in the kafka cluster
func (h *Writer) ProcessMessage(message, partition string) {
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka"})

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

// Shutdown cleans up the Kafka Writer
func (h *Writer) Shutdown() error {
	log := logging.GetLogger().WithFields(logrus.Fields{"package": "kafka"})
	log.Info("Shutting down Kafka Writer")

	// Defer close until function exists
	defer h.producer.Close()

	// Flush any remaining messages, waiting up until the flushTimeout
	msgs := h.producer.Flush(flushTimeout)
	if len(h.producer.ProduceChannel()) != 0 {
		return fmt.Errorf("%d messages were not flushed after a timeout of %d", msgs, flushTimeout)
	}

	return nil
}
