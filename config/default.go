package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GinMode string   `mapstructure:"GIN_MODE"`
	Env     string   `mapstructure:"ENV"`
	Port    string   `mapstructure:"PORT"`
	Origin  []string `mapstructure:"ORIGIN"`

	PgDbName    string `mapstructure:"POSTGRES_DB"`
	PgUser      string `mapstructure:"POSTGRES_USER"`
	PgPort      int    `mapstructure:"POSTGRES_PORT"`
	PgPassword  string `mapstructure:"POSTGRES_PASSWORD"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	PgLogLevel  string `mapstructure:"PGX_LOG_LEVEL"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	if os.Getenv("ENV") == "test" {
		viper.SetConfigName("env.test")
	} else {
		viper.SetConfigName(".env")
	}

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("failed read in config: %v", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
