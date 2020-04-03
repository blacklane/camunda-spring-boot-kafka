package config

import (
	"strings"
	"time"

	"github.com/caarlos0/env"
)

// Config - environment variables are parsed to this struct
type Config struct {
	AppName    string `env:"APP_NAME" envDefault:"worker"`
	Env        string `env:"ENV" envDefault:"development"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"debug"`
	LogOutput  string `env:"LOG_OUTPUT" envDefault:"console"`
	ServerPort int    `env:"PORT" envDefault:"8000"`

	RaygunAPIKey  string `env:"RAYGUN_API_KEY" envDefault:""`
	RaygunAppName string `env:"RAYGUN_APP_NAME" envDefault:"worker"`
	RaygunEnabled bool   `env:"RAYGUN_ENABLED" envDefault:"false"`

	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" envDefault:"5s"`
	// WriteTimeout maximum time the server will handle a request before timing out writes of the response.
	// It must be bigger than RequestTimeout
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"4s"`
	// RequestTimeout the timeout for the incoming request set on the request handler
	RequestTimeout    time.Duration `env:"REQUEST_TIMEOUT" envDefault:"2s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"1s"`

	// ShutdownTimeout the time the sever will wait server.Shutdown to return
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"6s"`

	// Database
	//DbHost     string `env:"DB_HOST" envDefault:""`
	//DbName     string `env:"DB_NAME" envDefault:""`
	//DbOptions  string `env:"DB_OPTIONS" envDefault:""`
	//DbPassword string `env:"DB_PASSWORD" envDefault:""`
	//DbPort     int    `env:"DB_PORT" envDefault:""`
	//DbUser     string `env:"DB_USER" envDefault:""`

	// Kafka
	KafkaBrokers       string `env:"KAFKA_BROKERS" envDefault:""`
	KafkaTopic         string `env:"KAFKA_TOPIC_RIDES" envDefault:""`
	KafkaConsumerGroup string `env:"KAFKA_CONSUMER_GROUP" envDefault:""`
}

//func (cfg *Config) PrepareDbURI() string {
//	return fmt.Sprintf(
//		"postgres://%s:%s@%s:%d/%s?%s",
//		cfg.DbUser,
//		cfg.DbPassword,
//		cfg.DbHost,
//		cfg.DbPort,
//		cfg.DbName,
//		cfg.DbOptions,
//	)
//}

func (cfg *Config) KafkaBrokersArray() []string {
	// convert comma separated list of brokers into string array
	return strings.Split(cfg.KafkaBrokers, ",")
}

// Parse environment variables, returns an error if an error occurs
func Parse() (Config, error) {
	confs := Config{}
	err := env.Parse(&confs)
	return confs, err
}
