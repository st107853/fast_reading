package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBname                 string        `mapstructure:"DBNAME"`
	DBuser                 string        `mapstructure:"DBUSER"`
	DBpassword             string        `mapstructure:"DBPASS"`
	Port                   string        `mapstructure:"PORT"`
	Host                   string        `mapstructure:"HOST"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	CORS1                  string        `mapstructure:"CORS_ORIGINS_1"`
	CORS2                  string        `mapstructure:"CORS_ORIGINS_2"`
	LogLevel               string        `mapstructure:"LOG_LEVEL"`  // "debug"|"info"|"warn"|"error"
	LogFormat              string        `mapstructure:"LOG_FORMAT"` // "json"|"text"
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
