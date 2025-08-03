package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBName               string        `mapstructure:"DB_NAME"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"ADDRESS"`
	SymmetricKey         string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_SESSION_DURATION"`
	MailgunDomain        string        `mapstructure:"EMAIL_DOMAIN"`
	MailgunAPIKEY        string        `mapstructure:"MAILGUN_API_KEY"`
	Domain               string        `mapstructure:"DOMAIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	viper.SetConfigName("secrets")
	_ = viper.MergeInConfig()

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	return config, nil
}
