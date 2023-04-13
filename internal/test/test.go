package test

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/fancar/tmp_xm/internal/config"
)

func init() {
	log.SetLevel(log.ErrorLevel)

}

// GetConfig returns the test configuration.
func GetConfig() config.Config {
	log.SetLevel(log.FatalLevel)

	var c config.Config

	c.PostgreSQL.DSN = "postgres://app_test@localhost:5442/app_test?sslmode=disable"
	c.PostgreSQL.Automigrate = false

	if v := os.Getenv("TEST_POSTGRES_DSN"); v != "" {
		c.PostgreSQL.DSN = v
	}

	c.Kafka.Brokers = []string{"172.30.0.1:9092"}
	c.Kafka.Topic = "epam-xm-test"
	c.Kafka.EventKeyTemplate = "company.{{ .Company }}.event.{{ .EventType }}"
	return c
}
