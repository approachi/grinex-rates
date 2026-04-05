package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GRPC     GRPCConfig
	Postgres PostgresConfig
	Metrics  MetricsConfig
	Tracing  TracingConfig
	Log      LogConfig
}

type GRPCConfig struct {
	Port int
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type MetricsConfig struct {
	Port int
}

type TracingConfig struct {
	Endpoint string
}

type LogConfig struct {
	Level string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetDefault("grpc.port", 50051)
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.database", "grinex")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("metrics.port", 8080)
	v.SetDefault("tracing.endpoint", "http://localhost:4318/v1/traces")
	v.SetDefault("log.level", "info")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
