/*
	Hans Nicolaus
	29 Aug 2023
*/

package model

type Config struct {
	GRPC  GRPC          `yaml:"grpc"`
	Kafka KafkaConsumer `yaml:"kafka_consumer"`
	Redis Redis         `yaml:"redis"`
}

type GRPC struct {
	Network string `yaml:"network"`
	Port    string `yaml:"port"`
}

type KafkaConsumer struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	GroupID string `yaml:"group_id"`
	Topic   string `yaml:"topic"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db" default:"0"`
}

var (
	DefaultConfigLocal Config = Config{
		GRPC: GRPC{
			Network: "tcp",
			Port:    ":50051",
		},
		Kafka: KafkaConsumer{
			Host:    "localhost",
			Port:    ":9092",
			GroupID: "bibit_consumer_group",
			Topic:   "bibit_challenge_1",
		},
		Redis: Redis{
			Host:     "localhost",
			Port:     ":6379",
			Password: "",
			DB:       0,
		},
	}
)
