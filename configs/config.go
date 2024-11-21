package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	WeatherAPIKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig(path string) (*conf, error) {
	cfg := &conf{}
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
