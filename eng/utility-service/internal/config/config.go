package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	Environment   string `mapstructure:"ENVIRONMENT"`
	AWSRegion     string `mapstructure:"AWS_REGION"`
	DocBucketName string `mapstructure:"DOC_BUCKET_NAME"`
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigFile("config.env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config.env file found, relying on System Env Variables")
		} else {
			log.Fatal("Error reading config file:", err)
		}
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal("Unable to decode into struct:", err)
	}

	if key := viper.GetString("AWS_ACCESS_KEY_ID"); key != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", key)
	}
	if secret := viper.GetString("AWS_SECRET_ACCESS_KEY"); secret != "" {
		os.Setenv("AWS_SECRET_ACCESS_KEY", secret)
	}
	
	if AppConfig.Port == "" {
		AppConfig.Port = "9004"
	}
	if AppConfig.AWSRegion == "" {
		AppConfig.AWSRegion = "ap-south-1"
	}
	if AppConfig.DocBucketName == "" {
		AppConfig.DocBucketName = "shikshakul-docs-dev"
	}
}
