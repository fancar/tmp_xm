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

	c.PostgreSQL.DSN = "postgres://app_test@localhost/app_test?sslmode=disable"
	//  c.PostgreSQL.DSN = "postgres://ernet_ns_test:ernet_ns_test@localhost/ernet_ns_test?sslmode=disable&binary_parameters=yes"
	c.PostgreSQL.Automigrate = false

	if v := os.Getenv("TEST_POSTGRES_DSN"); v != "" {
		c.PostgreSQL.DSN = v
	}

	c.CountryCheck.Enabled = true
	c.CountryCheck.UrlTmpl = "https://ipapi.co/{{ .IPaddress }}/country_name/"
	c.CountryCheck.CountryAllowed = "Cyprus"

	return c
}
