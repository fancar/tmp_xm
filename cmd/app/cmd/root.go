package cmd

import (
	"bytes"
	"io/ioutil"

	"github.com/fancar/tmp_xm/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appName = "epam_xm_exercise_v22"
	cfgFile string
	version string
)

// Execute executes the root command.
func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: "EPAM Golang Exercise v22",
	Long: `
	
	Mamaev Alexander fancatser@gmail.com
		
				- 2023 - 
	`,
	RunE: run,
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetDefault("postgre.max_open_connections", 10)

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config", "c", "", "path to configuration file (optional). Default config.toml")
	rootCmd.PersistentFlags().Int("log-level", 4, "debug=5, info=4, error=2, fatal=1, panic=0")

	viper.BindPFlag("general.log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	viper.SetDefault("postgre.dsn", "postgres://app@localhost/app?sslmode=disable")
	viper.SetDefault("postgre.max_idle_connections", 2)
	viper.SetDefault("postgre.max_open_connections", 10)
	viper.SetDefault("postgre.automigrate", true)

	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic", "epam-xm")
	viper.SetDefault("kafka.event_key_template", "company.{{ .Company }}.event.{{ .EventType }}")
	viper.SetDefault("kafka.mechanism", "PLAIN")
	viper.SetDefault("algorithm", "SHA512")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	config.Version = version
	config.AppName = appName

	if cfgFile != "" {
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
		viper.SetConfigType("toml")
		if err := viper.ReadConfig(bytes.NewBuffer(b)); err != nil {
			log.WithError(err).WithField("config", cfgFile).Fatal("error loading config file")
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/xm")
		if err := viper.ReadInConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
			default:
				log.WithError(err).Fatal("read configuration file error")
			}
		}
	}

	if err := viper.Unmarshal(&config.C); err != nil {
		log.WithError(err).Fatal("unmarshal config error")
	}
}
