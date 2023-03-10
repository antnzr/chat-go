package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GinMode     string `mapstructure:"GIN_MODE"`
	Port        string `mapstructure:"PORT"`
	Origin      string `mapstructure:"ORIGIN"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	AccessTokenSecret     string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge     int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenSecret    string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	RefreshTokenMaxAge    int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("failed read in config: %v", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
