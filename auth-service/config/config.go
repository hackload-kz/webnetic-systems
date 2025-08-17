package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     int    `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	RedisSSL			 bool   `mapstructure:"REDIS_SSL"`

	JWTSecret string `mapstructure:"JWT_SECRET"`
	JWTTTL    int    `mapstructure:"JWT_TTL"`

	AuthServicePort string `mapstructure:"AUTH_SERVICE_PORT"`

	PaymentEndpoint string `mapstructure:"PAYMENT_ENDPOINT"`
	EventProvider   string `mapstructure:"EVENT_PROVIDER"`
}

func New() (Config, error) {
	viper.AutomaticEnv() 
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("JWT_TTL", 3600)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal env: %w", err)
	}

	return cfg, nil
}
