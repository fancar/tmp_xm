package kafka

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"text/template"
	"time"

	"github.com/fancar/tmp_xm/internal/config"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	log "github.com/sirupsen/logrus"
)

var (
	wg               *sync.WaitGroup
	writer           *kafka.Writer // according to the documentation the Writer is thread safe
	eventKeyTemplate *template.Template
)

// Setup configures the kafka producer
func Setup(ctx context.Context, waitgroup *sync.WaitGroup, cfg config.Config) error {
	conf := cfg.Kafka
	wg = waitgroup
	if len(conf.Brokers) == 0 {
		log.Info("Kafka: no brokers specified. Skipped.")
		return nil
	}
	wc := kafka.WriterConfig{
		Async:    true,
		Brokers:  conf.Brokers,
		Topic:    conf.Topic,
		Balancer: &kafka.LeastBytes{},

		// Equal to kafka.DefaultDialer.
		// We do not want to use kafka.DefaultDialer itself, as we might modify
		// it below to setup SASLMechanism.
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
	}

	if conf.TLS {
		wc.Dialer.TLS = &tls.Config{}
	}

	if conf.Username != "" || conf.Password != "" {
		switch conf.Mechanism {
		case "PLAIN":
			wc.Dialer.SASLMechanism = plain.Mechanism{
				Username: conf.Username,
				Password: conf.Password,
			}
		case "SCRAM":
			var algorithm scram.Algorithm

			switch conf.Algorithm {
			case "SHA512":
				algorithm = scram.SHA512
			case "SHA256":
				algorithm = scram.SHA256
			default:
				return fmt.Errorf("unknown sasl algorithm %s", conf.Algorithm)
			}

			mechanism, err := scram.Mechanism(algorithm, conf.Username, conf.Password)
			if err != nil {
				return fmt.Errorf("sasl mechanism %w", err)
			}

			wc.Dialer.SASLMechanism = mechanism
		default:
			return fmt.Errorf("unknown sasl mechanism %s", conf.Mechanism)
		}
	}

	var err error
	eventKeyTemplate, err = template.New("key").Parse(conf.EventKeyTemplate)
	if err != nil {
		return fmt.Errorf("parse key template %w", err)
	}

	writer = kafka.NewWriter(wc)

	log.WithFields(log.Fields{
		"brokers":   conf.Brokers,
		"topic":     conf.Topic,
		"event_key": conf.EventKeyTemplate,
	}).Info("kafka: setup finished successfully!")

	return nil
}

// PublishMessage publishes the byte array recieved
func PublishMessage(ctx context.Context, company, event string, b []byte) error {
	wg.Add(1)
	defer wg.Done()

	keyBuf := bytes.NewBuffer(nil)

	err := eventKeyTemplate.Execute(keyBuf, struct {
		Company   string
		EventType string
	}{company, event})
	if err != nil {
		return fmt.Errorf("unable execute template %w", err)
	}

	key := keyBuf.Bytes()

	kmsg := kafka.Message{
		Value: b,
	}
	if len(key) > 0 {
		kmsg.Key = key
	}

	log.WithFields(log.Fields{
		"key":         string(key),
		"value_bytes": len(b),
		"value":       string(b),
		"company":     company,
		"event":       event,
	}).Info("kafka: message publishing ...")

	return writer.WriteMessages(context.Background(), kmsg)
}
