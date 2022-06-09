package cmd

import (
  "os"
  "text/template"

  "github.com/pkg/errors"
  "github.com/spf13/cobra"

  "github.com/fancar/tmp_xm/internal/config"
)

const configTemplate = `[general]
# This is an example of config.toml file for the program
#
# Log level
#
# debug=5, info=4, warning=3, error=2, fatal=1, panic=0
log_level={{ .General.LogLevel }}

# client's ip country check
# if enabled - only users from the country are allowed
# for Create\Delete operations
[country_check]
# true - to enable the feature
enabled="{{ .CountryCheck.Enabled }}"

# url template of ip-API service
url_tmpl="{{ .CountryCheck.UrlTmpl }}"

# short country name (eg "United States")
country_allowed="{{ .CountryCheck.CountryAllowed }}"

[external_api]
  # ip:port to bind the (user facing) http server to (web-interface and REST / gRPC api)
  bind="{{ .ExternalAPI.Bind }}"

  # http server TLS certificate (optional)
  tls_cert="{{ .ExternalAPI.TLSCert }}"

  # http server TLS key (optional)
  tls_key="{{ .ExternalAPI.TLSKey }}"

  # JWT secret used for api authentication / authorization
  # You could generate this by executing 'openssl rand -base64 32' for example
  jwt_secret="{{ .ExternalAPI.JWTSecret }}"

  # Allow origin header (CORS).
  #
  # Set this to allows cross-domain communication from the browser (CORS).
  # Example value: https://example.com.
  # When left blank (default), CORS will not be used.
  cors_allow_origin="{{ .ExternalAPI.CORSAllowOrigin }}"


# PostgreSQL settings.
#
[postgre]
  # PostgreSQL dsn (e.g.: postgres://user:password@hostname/database?sslmode=disable).
  #
  # Besides using an URL (e.g. 'postgres://user:password@hostname/database?sslmode=disable')
  # it is also possible to use the following format:
  # 'user=app dbname=app sslmode=disable'.
  #
  # The following connection parameters are supported:
  #
  # * dbname - The name of the database to connect to
  # * user - The user to sign in as
  # * password - The user's password
  # * host - The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)
  # * port - The port to bind to. (default is 5432)
  # * sslmode - Whether or not to use SSL (default is require, this is not the default for libpq)
  # * fallback_application_name - An application_name to fall back to if one isn't provided.
  # * connect_timeout - Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.
  # * sslcert - Cert file location. The file must contain PEM encoded data.
  # * sslkey - Key file location. The file must contain PEM encoded data.
  # * sslrootcert - The location of the root certificate file. The file must contain PEM encoded data.
  #
  # Valid values for sslmode are:
  #
  # * disable - No SSL
  # * require - Always SSL (skip verification)
  # * verify-ca - Always SSL (verify that the certificate presented by the server was signed by a trusted CA)
  # * verify-full - Always SSL (verify that the certification presented by the server was signed by a trusted CA and the server host name matches the one in the certificate)
  dsn="{{ .PostgreSQL.DSN }}"

  # Max open connections.
  #
  # This sets the max. number of open connections that are allowed in the
  # PostgreSQL connection pool (0 = unlimited).
  max_open_connections={{ .PostgreSQL.MaxOpenConnections }}

  # Max idle connections.
  #
  # This sets the max. number of idle connections in the PostgreSQL connection
  # pool (0 = no idle connections are retained).
  max_idle_connections={{ .PostgreSQL.MaxIdleConnections }}

`

var configCmd = &cobra.Command{
  Use:   "configfile",
  Short: "Print the configuration file",
  RunE: func(cmd *cobra.Command, args []string) error {
    t := template.Must(template.New("config").Parse(configTemplate))
    err := t.Execute(os.Stdout, &config.C)
    if err != nil {
      return errors.Wrap(err, "execute config template error")
    }
    return nil
  },
}
