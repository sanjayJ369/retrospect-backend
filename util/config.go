package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBName               string        `mapstructure:"DB_NAME"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	SymmetricKey         string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	MailgunDomain        string        `mapstructure:"MAILGUN_DOMAIN"`
	MailgunAPIKEY        string        `mapstructure:"MAILGUN_API_KEY"`
	Domain               string        `mapstructure:"DOMAIN"`
	TemplatesDir         string        `mapstructure:"TEMPLATES_DIR"`
	RedisURL             string        `mapstructure:"REDIS_URL"`
	RatelimitDuration    time.Duration `mapstructure:"RATELIMIT_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	viper.SetConfigName(".env")
	_ = viper.MergeInConfig()

	err = viper.Unmarshal(&config)
	return
}
