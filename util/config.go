package util

import "github.com/spf13/viper"

// Store all configuration of the app
// Read from the config file
type Config struct {
	DB_SOURCE string `mapstructure:"DB_SOURCE"`
	WEB_ADDR  string `mapstructure:"WEB_ADDR"`
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
