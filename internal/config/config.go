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

	CountryCheck struct {
		Enabled        bool   `mapstructure:"enabled"`
		UrlTmpl        string `mapstructure:"url_tmpl"`
		CountryAllowed string `mapstructure:"country_allowed"`
	} `mapstructure:"country_check"`

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
}

// C holds the global configuration.
var C Config
