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
	appName = "xm_exercise_v21"
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
	Short: "Golang Exercise - v21.0.0",
	Long: `
	Golang Exercise - v21.0.0
	By Mamaev Alexander fancatser@gmail.com
	via EPAM | 08.06.2022`,
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

	viper.SetDefault("country_check.enabled", true)
	viper.SetDefault("country_check.url_tmpl", "https://ipapi.co/{{ .IPaddress }}/country_name/")
	viper.SetDefault("country_check.country_allowed", "Cyprus")

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
