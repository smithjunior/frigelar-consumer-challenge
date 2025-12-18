package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka   KafkaConfig
	Service ServiceConfig
}

type KafkaConfig struct {
	BootstrapServers string `mapstructure:"KAFKA_BOOTSTRAP_SERVERS"`
	Username         string `mapstructure:"KAFKA_USERNAME"`
	Password         string `mapstructure:"KAFKA_PASSWORD"`
	SaslMechanism    string `mapstructure:"KAFKA_SASL_MECHANISM"`
	SecurityProtocol string `mapstructure:"KAFKA_SECURITY_PROTOCOL"`
	GroupId          string `mapstructure:"KAFKA_GROUP_ID"`
	Topic            string `mapstructure:"KAFKA_TOPIC"`
	TopicDLQ         string `mapstructure:"KAFKA_TOPIC_DLQ"`
}

type ServiceConfig struct {
	TargetServiceURL string `mapstructure:"TARGET_SERVICE_URL"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, relying on environment variables: %v", err)
	}

	var cfg Config
	
	cfg.Kafka.BootstrapServers = viper.GetString("KAFKA_BOOTSTRAP_SERVERS")
	cfg.Kafka.Username = viper.GetString("KAFKA_USERNAME")
	cfg.Kafka.Password = viper.GetString("KAFKA_PASSWORD")
	cfg.Kafka.SaslMechanism = viper.GetString("KAFKA_SASL_MECHANISM")
	cfg.Kafka.SecurityProtocol = viper.GetString("KAFKA_SECURITY_PROTOCOL")
	cfg.Kafka.GroupId = viper.GetString("KAFKA_GROUP_ID")
	cfg.Kafka.Topic = viper.GetString("KAFKA_TOPIC")
	cfg.Kafka.TopicDLQ = viper.GetString("KAFKA_TOPIC_DLQ")
	
	cfg.Service.TargetServiceURL = viper.GetString("TARGET_SERVICE_URL")

	return &cfg, nil
}
