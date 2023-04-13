package test

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/fancar/tmp_xm/internal/config"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

var (
	consumerCh chan *kafka.Message
)

// KafkaConsumer goroutine which reads remote topic.
func KafkaConsumer(conf config.Config) error {
	cfg := conf.Kafka
	f := log.Fields{
		"brokers": cfg.Brokers,
		"topic":   cfg.Topic,
	}

	consumerCh = make(chan *kafka.Message)

	if len(cfg.Brokers) == 0 {
		return fmt.Errorf("kafka-consumer: no brokers configured. Check you test config")
	}

	log.WithFields(f).Info("kafka-consumer: has been started")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		MinBytes: 1e2, // 1B
		MaxBytes: 1e6, // 1MB
	})

	go func() {
		defer r.Close()
		for {
			msg, err := r.ReadMessage(context.Background())
			if err != nil {
				//				return err
				if err == io.EOF {
					log.WithFields(f).Info("kafka-consumer: Reader goroutine stopped!")
					return
				}
				log.WithError(err).WithFields(f).Errorf(
					"kafka-consumer: unable to read messages. Sleep for 5s ...")
				time.Sleep(1 * time.Second)
				continue
			}
			log.WithFields(f).Debugf(
				"kafka-consumer: dl message recieved (%d bytes)", len(msg.Value))
			consumerCh <- &msg
		}
	}()
	return nil
}

// GetMessage wait for a message or timeout
func GetMessage(key string) (*kafka.Message, error) {
	for {
		select {
		case msg := <-consumerCh:
			if key == string(msg.Key) {
				return msg, nil
			}
		case <-time.After(time.Second):
			return nil, fmt.Errorf("timeout getting message from kafka")
		}
	}
}
