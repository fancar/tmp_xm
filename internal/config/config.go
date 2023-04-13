package config

import (
// "time"
)

// Version defines the version.
var (
	Version string
	AppName string
)

// Config defines the configuration.
type Config struct {
	General struct {
		LogLevel int `mapstructure:"log_level"`
	}

	ExternalAPI struct {
		Bind            string
		TLSCert         string `mapstructure:"tls_cert"`
		TLSKey          string `mapstructure:"tls_key"`
		JWTSecret       string `mapstructure:"jwt_secret"`
		CORSAllowOrigin string `mapstructure:"cors_allow_origin"`
	} `mapstructure:"external_api"`

	PostgreSQL struct {
		Automigrate        bool
		DSN                string `mapstructure:"dsn"`
		MaxOpenConnections int    `mapstructure:"max_open_connections"`
		MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	} `mapstructure:"postgre"`

	Kafka struct {
		// Owner            string                       `mapstructure:"owner"`
		Brokers          []string                     `mapstructure:"brokers"`
		TLS              bool                         `mapstructure:"tls"`
		Topic            string                       `mapstructure:"topic"` // main topic
		EventKeyTemplate string                       `mapstructure:"event_key_template"`
		Username         string                       `mapstructure:"username"`
		Password         string                       `mapstructure:"password"`
		Mechanism        string                       `mapstructure:"mechanism"`
		Algorithm        string                       `mapstructure:"algorithm"`
		Marshaler        string                       `mapstructure:"marshaler"`
		Reader           KafkaReaderConfig            `mapstructure:"reader"`
		Writers          map[string]KafkaWriterConfig `mapstructure:"writers"` // writer per event name
	}
}

// KafkaReaderConfig consumer cfg
type KafkaReaderConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Marshaler string `mapstructure:"marshaler"`
	Topic     string `mapstructure:"topic"`
	GroupID   string `mapstructure:"groupID"`
}

// KafkaWriterConfig producer event cfg
type KafkaWriterConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Marshaler string `mapstructure:"marshaler"`
	Topic     string `mapstructure:"topic"` // if empty - topic from root in use
	GroupID   string `mapstructure:"groupID"`
}

// C holds the global configuration.
var C Config
