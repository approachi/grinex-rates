package config

import (
	"fmt"

	"github.com/spf13/pflag"
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

	pflag.String("grpc-port", "", "gRPC server port")
	pflag.String("postgres-host", "", "PostgreSQL host")
	pflag.Int("postgres-port", 0, "PostgreSQL port")
	pflag.String("postgres-user", "", "PostgreSQL user")
	pflag.String("postgres-password", "", "PostgreSQL password")
	pflag.String("postgres-database", "", "PostgreSQL database name")
	pflag.String("postgres-sslmode", "", "PostgreSQL SSL mode (disable|require|verify-full)")
	pflag.Int("metrics-port", 0, "Prometheus metrics server port")
	pflag.String("tracing-endpoint", "", "OpenTelemetry collector endpoint")
	pflag.String("log-level", "", "Log level (debug, info, warn, error)")
	pflag.Parse()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Bind each flag explicitly to its viper key.
	// Priority: flag > config file
	for _, pair := range [][2]string{
		{"grpc.port", "grpc-port"},
		{"postgres.host", "postgres-host"},
		{"postgres.port", "postgres-port"},
		{"postgres.user", "postgres-user"},
		{"postgres.password", "postgres-password"},
		{"postgres.database", "postgres-database"},
		{"postgres.sslmode", "postgres-sslmode"},
		{"metrics.port", "metrics-port"},
		{"tracing.endpoint", "tracing-endpoint"},
		{"log.level", "log-level"},
	} {
		if err := v.BindPFlag(pair[0], pflag.Lookup(pair[1])); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := validate(&cfg); err != nil {
		panic(fmt.Sprintf("config: %s", err))
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	required := []struct {
		name string
		ok   bool
	}{
		{"grpc.port", cfg.GRPC.Port != 0},
		{"postgres.host", cfg.Postgres.Host != ""},
		{"postgres.port", cfg.Postgres.Port != 0},
		{"postgres.user", cfg.Postgres.User != ""},
		{"postgres.password", cfg.Postgres.Password != ""},
		{"postgres.database", cfg.Postgres.Database != ""},
		{"postgres.sslmode", cfg.Postgres.SSLMode != ""},
		{"metrics.port", cfg.Metrics.Port != 0},
		{"tracing.endpoint", cfg.Tracing.Endpoint != ""},
		{"log.level", cfg.Log.Level != ""},
	}
	for _, r := range required {
		if !r.ok {
			return fmt.Errorf("required key is missing or empty: %s", r.name)
		}
	}
	return nil
}
