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

	Redis struct {
		Servers    []string `mapstructure:"servers"`
		Cluster    bool     `mapstructure:"cluster"`
		MasterName string   `mapstructure:"master_name"`
		PoolSize   int      `mapstructure:"pool_size"`
		Password   string   `mapstructure:"password"`
		Database   int      `mapstructure:"database"`
	} `mapstructure:"redis"`

	PostgreSQL struct {
		DSN                string `mapstructure:"dsn"`
		MaxOpenConnections int    `mapstructure:"max_open_connections"`
		MaxIdleConnections int    `mapstructure:"max_idle_connections"`
	} `mapstructure:"postgre"`

	// Prometheus struct {
	// 	Bind string `mapstructure:"bind"`
	// } `mapstructure:"prometheus"`
}

// C holds the global configuration.
var C Config
