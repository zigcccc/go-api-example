package config

import (
	"log"

	"github.com/spf13/viper"
)

type ApiGoExampleConfig struct {
	Database struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		Port     int    `mapstructure:"port"`
		SslMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	OAuth struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		SecretKey    string `mapstructure:"secret_key"`
	} `mapstructure:"oauth"`
}

var apiGoExampleConfig ApiGoExampleConfig

func InitConfig() {
	viper.SetConfigName("api-go-example")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/api-go-example")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, using defaults")
	}

	if err := viper.Unmarshal(&apiGoExampleConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}

func GetConfig() *ApiGoExampleConfig {
	return &apiGoExampleConfig
}