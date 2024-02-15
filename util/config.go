package util

import (
	"time"

	"github.com/spf13/viper"
)

// Store all configuration of the app
// Read from the config file
type Config struct {
	ENVIROMENT             string        `mapstructure:"ENVIROMENT"`
	DB_SOURCE              string        `mapstructure:"DB_SOURCE"`
	WEB_ADDR               string        `mapstructure:"WEB_ADDR"`
	GRPC_ADDR              string        `mapstructure:"GRPC_ADDR"`
	MIGRATION_URL          string        `mapstructure:"MIGRATION_URL"`
	SYMMETRIC_KEY          string        `mapstructure:"SYMMETRIC_KEY"`
	ACCESS_DURATION        time.Duration `mapstructure:"ACCESS_DURATION"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	REDIS_ADDRESS          string        `mapstructure:"REDIS_ADDRESS"`
	EMAIL_SENDER_NAME      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EMAIL_SENDER_ADDRESS   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EMAIL_SENDER_PASSWORD  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig(path string) (con Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Auto reqrite envs when changed
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&con)
	return
}
